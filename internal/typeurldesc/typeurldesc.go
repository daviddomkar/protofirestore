package typeurldesc

import "google.golang.org/protobuf/reflect/protoreflect"

type TypeURLDescriptor struct {
	protoreflect.FieldDescriptor
}

func NewTypeURLDescriptor() TypeURLDescriptor {
	return TypeURLDescriptor{}
}

func (TypeURLDescriptor) Index() int {
	return -1
}

func (TypeURLDescriptor) FullName() protoreflect.FullName {
	return "@type"
}
