package protofirestore

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/DavidDomkar/protofirestore/internal/encoding/messageset"
	"github.com/DavidDomkar/protofirestore/internal/order"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func Marshal(m proto.Message) (map[string]interface{}, error) {
	return MarshalOptions{}.Marshal(m)
}

type MarshalOptions struct {
	// EmitUnpopulated specifies whether to emit unpopulated fields. It does not
	// emit unpopulated oneof fields or unpopulated extension fields.
	// The JSON value emitted for unpopulated fields are as follows:
	//  ╔═══════╤════════════════════════════╗
	//  ║ JSON  │ Protobuf field             ║
	//  ╠═══════╪════════════════════════════╣
	//  ║ false │ proto3 boolean fields      ║
	//  ║ 0     │ proto3 numeric fields      ║
	//  ║ ""    │ proto3 string/bytes fields ║
	//  ║ null  │ proto2 scalar fields       ║
	//  ║ null  │ message fields             ║
	//  ║ []    │ list fields                ║
	//  ║ {}    │ map fields                 ║
	//  ╚═══════╧════════════════════════════╝
	EmitUnpopulated bool

	// EmitDefaultValues specifies whether to emit default-valued primitive fields,
	// empty lists, and empty maps. The fields affected are as follows:
	//  ╔═══════╤════════════════════════════════════════╗
	//  ║ JSON  │ Protobuf field                         ║
	//  ╠═══════╪════════════════════════════════════════╣
	//  ║ false │ non-optional scalar boolean fields     ║
	//  ║ 0     │ non-optional scalar numeric fields     ║
	//  ║ ""    │ non-optional scalar string/byte fields ║
	//  ║ []    │ empty repeated fields                  ║
	//  ║ {}    │ empty map fields                       ║
	//  ╚═══════╧════════════════════════════════════════╝
	//
	// Behaves similarly to EmitUnpopulated, but does not emit "null"-value fields,
	// i.e. presence-sensing fields that are omitted will remain omitted to preserve
	// presence-sensing.
	// EmitUnpopulated takes precedence over EmitDefaultValues since the former generates
	// a strict superset of the latter.
	EmitDefaultValues bool

	// Resolver is used for looking up types when expanding google.protobuf.Any
	// messages. If nil, this defaults to using protoregistry.GlobalTypes.
	Resolver interface {
		protoregistry.ExtensionTypeResolver
		protoregistry.MessageTypeResolver
	}
}

func (o MarshalOptions) Marshal(m proto.Message) (map[string]interface{}, error) {
	return o.marshal(m)
}

func (o MarshalOptions) marshal(m proto.Message) (map[string]interface{}, error) {
	if o.Resolver == nil {
		o.Resolver = protoregistry.GlobalTypes
	}

	if m == nil {
		return make(map[string]interface{}), nil
	}

	enc := encoder{o}

	if marshal := wellKnownTypeMarshaler(m.ProtoReflect().Descriptor().FullName()); marshal != nil {
		return nil, errors.New("no support for well known types as top level objects in firestore documents")
	}

	if object, err := enc.marshalMessage(m.ProtoReflect()); err != nil {
		return nil, err
	} else {
		return object, proto.CheckInitialized(m)
	}
}

type encoder struct {
	opts MarshalOptions
}

// unpopulatedFieldRanger wraps a protoreflect.Message and modifies its Range
// method to additionally iterate over unpopulated fields.
type unpopulatedFieldRanger struct {
	protoreflect.Message

	skipNull bool
}

func (m unpopulatedFieldRanger) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	fds := m.Descriptor().Fields()
	for i := 0; i < fds.Len(); i++ {
		fd := fds.Get(i)
		if m.Has(fd) || fd.ContainingOneof() != nil {
			continue // ignore populated fields and fields within a oneofs
		}

		v := m.Get(fd)
		isProto2Scalar := fd.Syntax() == protoreflect.Proto2 && fd.Default().IsValid()
		isSingularMessage := fd.Cardinality() != protoreflect.Repeated && fd.Message() != nil
		if isProto2Scalar || isSingularMessage {
			if m.skipNull {
				continue
			}
			v = protoreflect.Value{} // use invalid value to emit null
		}
		if !f(fd, v) {
			return
		}
	}
	m.Message.Range(f)
}

// marshalMessage marshals the fields in the given protoreflect.Message.
// If the typeURL is non-empty, then a synthetic "@type" field is injected
// containing the URL as the value.
func (e encoder) marshalMessage(m protoreflect.Message) (map[string]interface{}, error) {
	if messageset.IsMessageSet(m.Descriptor()) {
		return nil, errors.New("no support for proto1 MessageSets")
	}

	var fields order.FieldRanger = m
	switch {
	case e.opts.EmitUnpopulated:
		fields = unpopulatedFieldRanger{Message: m, skipNull: false}
	case e.opts.EmitDefaultValues:
		fields = unpopulatedFieldRanger{Message: m, skipNull: true}
	}

	object := make(map[string]interface{})

	var err error
	order.RangeFields(fields, order.IndexNameFieldOrder, func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if value, e := e.marshalValue(v, fd); e != nil {
			err = e
			return false
		} else if value != nil {
			name := fd.JSONName()
			object[name] = value
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	return object, nil
}

// marshalValue marshals the given protoreflect.Value.
func (e encoder) marshalValue(val protoreflect.Value, fd protoreflect.FieldDescriptor) (interface{}, error) {
	switch {
	case fd.IsList():
		return e.marshalList(val.List(), fd)
	case fd.IsMap():
		mmap, err := e.marshalMap(val.Map(), fd)

		if err != nil {
			return nil, err
		}

		if len(mmap) == 0 {
			return nil, nil
		}

		return mmap, nil
	default:
		return e.marshalSingular(val, fd)
	}
}

func (e encoder) marshalSingular(val protoreflect.Value, fd protoreflect.FieldDescriptor) (interface{}, error) {
	if !val.IsValid() {
		return nil, nil
	}

	switch kind := fd.Kind(); kind {
	case protoreflect.BoolKind:
		return val.Bool(), nil

	case protoreflect.StringKind:
		if val.String() == "" {
			return nil, nil
		}

		if !utf8.ValidString(val.String()) {
			return nil, fmt.Errorf("field %v contains invalid UTF-8", string(fd.FullName()))
		}

		return val.String(), nil

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return int32(val.Int()), nil

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return val.Int(), nil

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return uint32(val.Uint()), nil

	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return val.Uint(), nil

	case protoreflect.FloatKind:
		return float32(val.Float()), nil

	case protoreflect.DoubleKind:
		return val.Float(), nil

	case protoreflect.BytesKind:
		if len(val.Bytes()) == 0 {
			return nil, nil
		}
		return val.Bytes(), nil

	case protoreflect.EnumKind:
		if fd.Enum().FullName() == "google.protobuf.NullValue" {
			return nil, nil
		} else {
			desc := fd.Enum().Values().ByNumber(val.Enum())

			if desc == nil {
				return int64(val.Enum()), nil
			}

			return string(desc.Name()), nil
		}

	case protoreflect.MessageKind, protoreflect.GroupKind:
		m := val.Message()

		if marshal := wellKnownTypeMarshaler(m.Descriptor().FullName()); marshal != nil {
			return marshal(e, m)
		}

		if object, err := e.marshalMessage(m); err != nil {
			return nil, err
		} else if len(object) != 0 {
			return object, nil
		}

		return nil, nil

	default:
		panic(fmt.Sprintf("%v has unknown kind: %v", fd.FullName(), kind))
	}
}

// marshalList marshals the given protoreflect.List.
func (e encoder) marshalList(list protoreflect.List, fd protoreflect.FieldDescriptor) ([]interface{}, error) {
	if list.Len() == 0 {
		return nil, nil
	}

	array := make([]interface{}, list.Len())

	for i := 0; i < list.Len(); i++ {
		item := list.Get(i)
		if value, err := e.marshalSingular(item, fd); err != nil {
			return nil, err
		} else if value != nil {
			array[i] = value
		}
	}

	if len(array) == 0 {
		return nil, nil
	}

	return array, nil
}

// marshalMap marshals given protoreflect.Map.
func (e encoder) marshalMap(mmap protoreflect.Map, fd protoreflect.FieldDescriptor) (map[string]interface{}, error) {
	if mmap.Len() == 0 {
		return nil, nil
	}

	object := make(map[string]interface{})

	var err error
	order.RangeEntries(mmap, order.GenericKeyOrder, func(k protoreflect.MapKey, v protoreflect.Value) bool {
		if value, e := e.marshalSingular(v, fd.MapValue()); e != nil {
			err = e
			return false
		} else if value != nil {
			name := k.String()
			object[name] = value
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	return object, nil
}
