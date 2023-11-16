package pkg_test

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/proto"

	pb3 "github.com/DavidDomkar/protofirestore/internal/testprotos/textpb3"
	pkg "github.com/DavidDomkar/protofirestore/pkg"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		desc  string
		input proto.Message
		want  map[string]interface{}
	}{
		{
			desc:  "proto3 scalars not set",
			input: &pb3.Scalars{},
			want:  map[string]interface{}{},
		},
		{
			desc:  "proto3 optional not set",
			input: &pb3.Proto3Optional{},
			want:  map[string]interface{}{},
		},
		{
			desc: "proto3 optional set to zero values",
			input: &pb3.Proto3Optional{
				OptBool:    proto.Bool(false),
				OptInt32:   proto.Int32(0),
				OptInt64:   proto.Int64(0),
				OptUint32:  proto.Uint32(0),
				OptUint64:  proto.Uint64(0),
				OptFloat:   proto.Float32(0),
				OptDouble:  proto.Float64(0),
				OptString:  proto.String(""),
				OptBytes:   []byte{},
				OptEnum:    pb3.Enum_ZERO.Enum(),
				OptMessage: &pb3.Nested{},
			},
			want: map[string]interface{}{
				"optBool":   false,
				"optInt32":  0,
				"optInt64":  0,
				"optUint32": 0,
				"optUint64": 0,
				"optFloat":  0,
				"optDouble": 0,
				"optEnum":   "ZERO",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {

			got, err := pkg.Marshal(tt.input)

			if err != nil {
				t.Errorf("Marshal() error = %v", err)
				return
			}

			equal := reflect.DeepEqual(got, tt.want)

			if !equal {
				t.Errorf("Marshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
