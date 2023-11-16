package pkg

import "google.golang.org/protobuf/reflect/protoreflect"

type marshalFunc func(encoder, protoreflect.Message) (map[string]interface{}, error)

func wellKnownTypeMarshaler(name protoreflect.FullName) marshalFunc {
	// TODO: implement well known types

	return nil
}
