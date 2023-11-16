// Package flags provides a set of flags controlled by build tags.
package flags

// ProtoLegacy specifies whether to enable support for legacy functionality
// such as MessageSets, weak fields, and various other obscure behavior
// that is necessary to maintain backwards compatibility with proto1 or
// the pre-release variants of proto2 and proto3.
//
// This is disabled by default unless built with the "protolegacy" tag.
//
// WARNING: The compatibility agreement covers nothing provided by this flag.
// As such, functionality may suddenly be removed or changed at our discretion.
const ProtoLegacy = protoLegacy