package pkg_test

import (
	"math"
	"testing"

	"google.golang.org/protobuf/proto"

	pb2 "github.com/DavidDomkar/protofirestore/internal/testprotos/textpb2"
	pb3 "github.com/DavidDomkar/protofirestore/internal/testprotos/textpb3"
	pkg "github.com/DavidDomkar/protofirestore/pkg"
	"github.com/go-test/deep"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		desc    string
		input   proto.Message
		want    map[string]interface{}
		wantErr bool
	}{
		{
			desc:  "proto2 optional scalars not set",
			input: &pb2.Scalars{},
			want:  map[string]interface{}{},
		}, {
			desc:  "proto3 scalars not set",
			input: &pb3.Scalars{},
			want:  map[string]interface{}{},
		}, {
			desc:  "proto3 optional not set",
			input: &pb3.Proto3Optional{},
			want:  map[string]interface{}{},
		}, {
			desc: "proto2 optional scalars set to zero values",
			input: &pb2.Scalars{
				OptBool:     proto.Bool(false),
				OptInt32:    proto.Int32(0),
				OptInt64:    proto.Int64(0),
				OptUint32:   proto.Uint32(0),
				OptUint64:   proto.Uint64(0),
				OptSint32:   proto.Int32(0),
				OptSint64:   proto.Int64(0),
				OptFixed32:  proto.Uint32(0),
				OptFixed64:  proto.Uint64(0),
				OptSfixed32: proto.Int32(0),
				OptSfixed64: proto.Int64(0),
				OptFloat:    proto.Float32(0),
				OptDouble:   proto.Float64(0),
				OptBytes:    []byte{},
				OptString:   proto.String(""),
			},
			want: map[string]interface{}{
				"optBool":     false,
				"optInt32":    int32(0),
				"optInt64":    int64(0),
				"optUint32":   uint32(0),
				"optUint64":   uint64(0),
				"optSint32":   int32(0),
				"optSint64":   int64(0),
				"optFixed32":  uint32(0),
				"optFixed64":  uint64(0),
				"optSfixed32": int32(0),
				"optSfixed64": int64(0),
				"optFloat":    float32(0),
				"optDouble":   float64(0),
			},
		}, {
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
				"optInt32":  int32(0),
				"optInt64":  int64(0),
				"optUint32": uint32(0),
				"optUint64": uint64(0),
				"optFloat":  float32(0),
				"optDouble": float64(0),
				"optEnum":   "ZERO",
			},
		}, {
			desc: "proto2 optional scalars set to some values",
			input: &pb2.Scalars{
				OptBool:     proto.Bool(true),
				OptInt32:    proto.Int32(0xff),
				OptInt64:    proto.Int64(0xdeadbeef),
				OptUint32:   proto.Uint32(47),
				OptUint64:   proto.Uint64(0xdeadbeef),
				OptSint32:   proto.Int32(-1001),
				OptSint64:   proto.Int64(-0xffff),
				OptFixed64:  proto.Uint64(64),
				OptSfixed32: proto.Int32(-32),
				OptFloat:    proto.Float32(1.02),
				OptDouble:   proto.Float64(1.234),
				OptBytes:    []byte("谷歌"),
				OptString:   proto.String("谷歌"),
			},
			want: map[string]interface{}{
				"optBool":     true,
				"optInt32":    int32(255),
				"optInt64":    int64(3735928559),
				"optUint32":   uint32(47),
				"optUint64":   uint64(3735928559),
				"optSint32":   int32(-1001),
				"optSint64":   int64(-65535),
				"optFixed64":  uint64(64),
				"optSfixed32": int32(-32),
				"optFloat":    float32(1.02),
				"optDouble":   float64(1.234),
				"optBytes":    []byte("谷歌"),
				"optString":   "谷歌",
			},
		}, {
			desc: "string",
			input: &pb3.Scalars{
				SString: "谷歌",
			},
			want: map[string]interface{}{
				"sString": "谷歌",
			},
		}, {
			desc: "string with invalid UTF8",
			input: &pb3.Scalars{
				SString: "abc\xff",
			},
			wantErr: true,
		}, {
			desc: "float nan",
			input: &pb3.Scalars{
				SFloat: float32(math.NaN()),
			},
			want: map[string]interface{}{
				"sFloat": float32(math.NaN()),
			},
		}, {
			desc: "float positive infinity",
			input: &pb3.Scalars{
				SFloat: float32(math.Inf(1)),
			},
			want: map[string]interface{}{
				"sFloat": float32(math.Inf(1)),
			},
		}, {
			desc: "float negative infinity",
			input: &pb3.Scalars{
				SFloat: float32(math.Inf(-1)),
			},
			want: map[string]interface{}{
				"sFloat": float32(math.Inf(-1)),
			},
		}, {
			desc: "double nan",
			input: &pb3.Scalars{
				SDouble: math.NaN(),
			},
			want: map[string]interface{}{
				"sDouble": math.NaN(),
			},
		}, {
			desc: "double positive infinity",
			input: &pb3.Scalars{
				SDouble: math.Inf(1),
			},
			want: map[string]interface{}{
				"sDouble": math.Inf(1),
			},
		}, {
			desc: "double negative infinity",
			input: &pb3.Scalars{
				SDouble: math.Inf(-1),
			},
			want: map[string]interface{}{
				"sDouble": math.Inf(-1),
			},
		}, {
			desc:  "proto2 enum not set",
			input: &pb2.Enums{},
			want:  map[string]interface{}{},
		}, {
			desc: "proto2 enum set to zero value",
			input: &pb2.Enums{
				OptEnum:       pb2.Enum(0).Enum(),
				OptNestedEnum: pb2.Enums_NestedEnum(0).Enum(),
			},
			want: map[string]interface{}{
				"optEnum":       int64(0),
				"optNestedEnum": int64(0),
			},
		}, {
			desc: "proto2 enum",
			input: &pb2.Enums{
				OptEnum:       pb2.Enum_ONE.Enum(),
				OptNestedEnum: pb2.Enums_UNO.Enum(),
			},
			want: map[string]interface{}{
				"optEnum":       "ONE",
				"optNestedEnum": "UNO",
			},
		}, {
			desc: "proto2 enum set to numeric values",
			input: &pb2.Enums{
				OptEnum:       pb2.Enum(2).Enum(),
				OptNestedEnum: pb2.Enums_NestedEnum(2).Enum(),
			},
			want: map[string]interface{}{
				"optEnum":       "TWO",
				"optNestedEnum": "DOS",
			},
		}, {
			desc: "proto2 enum set to unnamed numeric values",
			input: &pb2.Enums{
				OptEnum:       pb2.Enum(101).Enum(),
				OptNestedEnum: pb2.Enums_NestedEnum(-101).Enum(),
			},
			want: map[string]interface{}{
				"optEnum":       int64(101),
				"optNestedEnum": int64(-101),
			},
		}, {
			desc:  "proto3 enum not set",
			input: &pb3.Enums{},
			want:  map[string]interface{}{},
		}, {
			desc: "proto3 enum set to zero value",
			input: &pb3.Enums{
				SEnum:       pb3.Enum_ZERO,
				SNestedEnum: pb3.Enums_CERO,
			},
			want: map[string]interface{}{},
		}, {
			desc: "proto3 enum",
			input: &pb3.Enums{
				SEnum:       pb3.Enum_ONE,
				SNestedEnum: pb3.Enums_UNO,
			},
			want: map[string]interface{}{
				"sEnum":       "ONE",
				"sNestedEnum": "UNO",
			},
		}, {
			desc: "proto3 enum set to numeric values",
			input: &pb3.Enums{
				SEnum:       2,
				SNestedEnum: 2,
			},
			want: map[string]interface{}{
				"sEnum":       "TWO",
				"sNestedEnum": "DOS",
			},
		}, {
			desc: "proto3 enum set to unnamed numeric values",
			input: &pb3.Enums{
				SEnum:       -47,
				SNestedEnum: 47,
			},
			want: map[string]interface{}{
				"sEnum":       int64(-47),
				"sNestedEnum": int64(47),
			},
		}, {
			desc:  "proto2 nested message not set",
			input: &pb2.Nests{},
			want:  map[string]interface{}{},
		}, {
			desc: "proto2 nested message set to empty",
			input: &pb2.Nests{
				OptNested: &pb2.Nested{},
				Optgroup:  &pb2.Nests_OptGroup{},
			},
			want: map[string]interface{}{},
		}, {
			desc: "proto2 nested messages",
			input: &pb2.Nests{
				OptNested: &pb2.Nested{
					OptString: proto.String("nested message"),
					OptNested: &pb2.Nested{
						OptString: proto.String("another nested message"),
					},
				},
			},
			want: map[string]interface{}{
				"optNested": map[string]interface{}{
					"optString": "nested message",
					"optNested": map[string]interface{}{
						"optString": "another nested message",
					},
				},
			},
		}, {
			desc: "proto2 groups",
			input: &pb2.Nests{
				Optgroup: &pb2.Nests_OptGroup{
					OptString: proto.String("inside a group"),
					OptNested: &pb2.Nested{
						OptString: proto.String("nested message inside a group"),
					},
					Optnestedgroup: &pb2.Nests_OptGroup_OptNestedGroup{
						OptFixed32: proto.Uint32(47),
					},
				},
			},
			want: map[string]interface{}{
				"optgroup": map[string]interface{}{
					"optString": "inside a group",
					"optNested": map[string]interface{}{
						"optString": "nested message inside a group",
					},
					"optnestedgroup": map[string]interface{}{
						"optFixed32": uint32(47),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {

			got, err := pkg.Marshal(tt.input)

			if err != nil && !tt.wantErr {
				t.Errorf("Marshal() returned error: %v\n", err)
			}

			if err == nil && tt.wantErr {
				t.Errorf("Marshal() got nil error, want error\n")
			}

			diff := deep.Equal(got, tt.want)

			if diff != nil {
				t.Error(diff)
			}
		})
	}
}
