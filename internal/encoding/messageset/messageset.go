package messageset

import "google.golang.org/protobuf/reflect/protoreflect"

// IsMessageSet returns whether the message uses the MessageSet wire format.
func IsMessageSet(md protoreflect.MessageDescriptor) bool {
	xmd, ok := md.(interface{ IsMessageSet() bool })
	return ok && xmd.IsMessageSet()
}
