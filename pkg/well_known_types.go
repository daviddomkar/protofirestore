package pkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/DavidDomkar/protofirestore/internal/genid"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type marshalFunc func(encoder, protoreflect.Message) (interface{}, error)

func wellKnownTypeMarshaler(name protoreflect.FullName) marshalFunc {
	if name.Parent() == genid.GoogleProtobuf_package {
		switch name.Name() {
		case genid.Any_message_name:
			return encoder.marshalAny
		case genid.Timestamp_message_name:
			return encoder.marshalTimestamp
		case genid.Duration_message_name:
			return encoder.marshalDuration
		case genid.BoolValue_message_name,
			genid.Int32Value_message_name,
			genid.Int64Value_message_name,
			genid.UInt32Value_message_name,
			genid.UInt64Value_message_name,
			genid.FloatValue_message_name,
			genid.DoubleValue_message_name,
			genid.StringValue_message_name,
			genid.BytesValue_message_name:
			return encoder.marshalWrapperType
		case genid.Struct_message_name:
			return encoder.marshalStruct
		case genid.ListValue_message_name:
			return encoder.marshalListValue
		case genid.Value_message_name:
			return encoder.marshalKnownValue
		case genid.FieldMask_message_name:
			return encoder.marshalFieldMask
		case genid.Empty_message_name:
			return encoder.marshalEmpty
		}
	}
	return nil
}

func (e encoder) marshalAny(m protoreflect.Message) (interface{}, error) {
	return nil, errors.New("no support for any well known type")
}

const (
	maxTimestampSeconds = 253402300799
	minTimestampSeconds = -62135596800
)

func (e encoder) marshalTimestamp(m protoreflect.Message) (interface{}, error) {
	fds := m.Descriptor().Fields()
	fdSeconds := fds.ByNumber(genid.Timestamp_Seconds_field_number)
	fdNanos := fds.ByNumber(genid.Timestamp_Nanos_field_number)

	secsVal := m.Get(fdSeconds)
	nanosVal := m.Get(fdNanos)
	secs := secsVal.Int()
	nanos := nanosVal.Int()

	if secs < minTimestampSeconds || secs > maxTimestampSeconds {
		return nil, fmt.Errorf("%s: seconds out of range %v", genid.Timestamp_message_fullname, secs)
	}

	if nanos < 0 || nanos > secondsInNanos {
		return nil, fmt.Errorf("%s: nanos out of range %v", genid.Timestamp_message_fullname, nanos)
	}

	return time.Unix(secs, nanos).UTC(), nil
}

const (
	secondsInNanos       = 999999999
	maxSecondsInDuration = 315576000000
)

func (e encoder) marshalDuration(m protoreflect.Message) (interface{}, error) {
	return nil, errors.New("no support for duration well known type")
}

func (e encoder) marshalWrapperType(m protoreflect.Message) (interface{}, error) {
	return nil, errors.New("no support for wrapper type well known type")
}

func (e encoder) marshalStruct(m protoreflect.Message) (interface{}, error) {
	return nil, errors.New("no support for struct well known type")
}

func (e encoder) marshalListValue(m protoreflect.Message) (interface{}, error) {
	return nil, errors.New("no support for list value well known type")
}

func (e encoder) marshalKnownValue(m protoreflect.Message) (interface{}, error) {
	return nil, errors.New("no support for value well known type")
}

func (e encoder) marshalFieldMask(m protoreflect.Message) (interface{}, error) {
	return nil, errors.New("no support for field mask well known type")
}

func (e encoder) marshalEmpty(m protoreflect.Message) (interface{}, error) {
	return nil, nil
}
