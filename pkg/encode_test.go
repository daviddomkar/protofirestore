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
		}, {
			desc:  "proto3 nested message not set",
			input: &pb3.Nests{},
			want:  map[string]interface{}{},
		}, {
			desc: "proto3 nested message set to empty",
			input: &pb3.Nests{
				SNested: &pb3.Nested{},
			},
			want: map[string]interface{}{},
		}, {
			desc: "proto3 nested message",
			input: &pb3.Nests{
				SNested: &pb3.Nested{
					SString: "nested message",
					SNested: &pb3.Nested{
						SString: "another nested message",
					},
				},
			},
			want: map[string]interface{}{
				"sNested": map[string]interface{}{
					"sString": "nested message",
					"sNested": map[string]interface{}{
						"sString": "another nested message",
					},
				},
			},
		}, {
			desc:  "oneof not set",
			input: &pb3.Oneofs{},
			want:  map[string]interface{}{},
		}, {
			desc: "oneof set to empty string",
			input: &pb3.Oneofs{
				Union: &pb3.Oneofs_OneofString{},
			},
			want: map[string]interface{}{},
		}, {
			desc: "oneof set to string",
			input: &pb3.Oneofs{
				Union: &pb3.Oneofs_OneofString{
					OneofString: "hello",
				},
			},
			want: map[string]interface{}{
				"oneofString": "hello",
			},
		}, {
			desc: "oneof set to enum",
			input: &pb3.Oneofs{
				Union: &pb3.Oneofs_OneofEnum{
					OneofEnum: pb3.Enum_ZERO,
				},
			},
			want: map[string]interface{}{
				"oneofEnum": "ZERO",
			},
		}, {
			desc: "oneof set to empty message",
			input: &pb3.Oneofs{
				Union: &pb3.Oneofs_OneofNested{
					OneofNested: &pb3.Nested{},
				},
			},
			want: map[string]interface{}{},
		}, {
			desc: "oneof set to message",
			input: &pb3.Oneofs{
				Union: &pb3.Oneofs_OneofNested{
					OneofNested: &pb3.Nested{
						SString: "nested message",
					},
				},
			},
			want: map[string]interface{}{
				"oneofNested": map[string]interface{}{
					"sString": "nested message",
				},
			},
		}, {
			desc:  "repeated fields not set",
			input: &pb2.Repeats{},
			want:  map[string]interface{}{},
		}, {
			desc: "repeated fields set to empty slices",
			input: &pb2.Repeats{
				RptBool:   []bool{},
				RptInt32:  []int32{},
				RptInt64:  []int64{},
				RptUint32: []uint32{},
				RptUint64: []uint64{},
				RptFloat:  []float32{},
				RptDouble: []float64{},
				RptBytes:  [][]byte{},
			},
			want: map[string]interface{}{},
		}, {
			desc: "repeated fields set to some values",
			input: &pb2.Repeats{
				RptBool:   []bool{true, false, true, true},
				RptInt32:  []int32{1, 6, 0, 0},
				RptInt64:  []int64{-64, 47},
				RptUint32: []uint32{0xff, 0xffff},
				RptUint64: []uint64{0xdeadbeef},
				RptFloat:  []float32{float32(math.NaN()), float32(math.Inf(1)), float32(math.Inf(-1)), 1.034},
				RptDouble: []float64{math.NaN(), math.Inf(1), math.Inf(-1), 1.23e-308},
				RptString: []string{"hello", "世界"},
				RptBytes: [][]byte{
					[]byte("hello"),
					[]byte("\xe4\xb8\x96\xe7\x95\x8c"),
				},
			},
			want: map[string]interface{}{
				"rptBool":   []interface{}{true, false, true, true},
				"rptInt32":  []interface{}{int32(1), int32(6), int32(0), int32(0)},
				"rptInt64":  []interface{}{int64(-64), int64(47)},
				"rptUint32": []interface{}{uint32(255), uint32(65535)},
				"rptUint64": []interface{}{uint64(3735928559)},
				"rptFloat":  []interface{}{float32(math.NaN()), float32(math.Inf(1)), float32(math.Inf(-1)), float32(1.034)},
				"rptDouble": []interface{}{math.NaN(), math.Inf(1), math.Inf(-1), 1.23e-308},
				"rptString": []interface{}{"hello", "世界"},
				"rptBytes": []interface{}{
					[]byte("hello"),
					[]byte("\xe4\xb8\x96\xe7\x95\x8c"),
				},
			},
		}, {
			desc: "repeated enums",
			input: &pb2.Enums{
				RptEnum:       []pb2.Enum{pb2.Enum_ONE, 2, pb2.Enum_TEN, 42},
				RptNestedEnum: []pb2.Enums_NestedEnum{2, 47, 10},
			},
			want: map[string]interface{}{
				"rptEnum":       []interface{}{"ONE", "TWO", "TEN", int64(42)},
				"rptNestedEnum": []interface{}{"DOS", int64(47), "DIEZ"},
			},
		}, {
			desc: "repeated messages set to empty",
			input: &pb2.Nests{
				RptNested: []*pb2.Nested{},
				Rptgroup:  []*pb2.Nests_RptGroup{},
			},
			want: map[string]interface{}{},
		}, {
			desc: "repeated messages",
			input: &pb2.Nests{
				RptNested: []*pb2.Nested{
					{
						OptString: proto.String("repeat nested one"),
					},
					{
						OptString: proto.String("repeat nested two"),
						OptNested: &pb2.Nested{
							OptString: proto.String("inside repeat nested two"),
						},
					},
					{},
				},
			},
			want: map[string]interface{}{
				"rptNested": []interface{}{
					map[string]interface{}{
						"optString": "repeat nested one",
					},
					map[string]interface{}{
						"optString": "repeat nested two",
						"optNested": map[string]interface{}{
							"optString": "inside repeat nested two",
						},
					},
					nil,
				},
			},
		}, {
			desc: "repeated messages contains nil value",
			input: &pb2.Nests{
				RptNested: []*pb2.Nested{nil, {}},
			},
			want: map[string]interface{}{
				"rptNested": []interface{}{nil, nil},
			},
		}, {
			desc: "repeated groups",
			input: &pb2.Nests{
				Rptgroup: []*pb2.Nests_RptGroup{
					{
						RptString: []string{"hello", "world"},
					},
					{},
					nil,
				},
			},
			want: map[string]interface{}{
				"rptgroup": []interface{}{
					map[string]interface{}{
						"rptString": []interface{}{"hello", "world"},
					},
					nil,
					nil,
				},
			},
		}, {
			desc:  "map fields not set",
			input: &pb3.Maps{},
			want:  map[string]interface{}{},
		}, {
			desc: "map fields set to empty",
			input: &pb3.Maps{
				Int32ToStr:   map[int32]string{},
				BoolToUint32: map[bool]uint32{},
				Uint64ToEnum: map[uint64]pb3.Enum{},
				StrToNested:  map[string]*pb3.Nested{},
				StrToOneofs:  map[string]*pb3.Oneofs{},
			},
			want: map[string]interface{}{},
		}, {
			desc: "map fields 1",
			input: &pb3.Maps{
				BoolToUint32: map[bool]uint32{
					true:  42,
					false: 101,
				},
			},
			want: map[string]interface{}{
				"boolToUint32": map[string]interface{}{
					"false": uint32(101),
					"true":  uint32(42),
				},
			},
		}, {
			desc: "map fields 2",
			input: &pb3.Maps{
				Int32ToStr: map[int32]string{
					-101: "-101",
					0xff: "0xff",
					0:    "zero",
				},
			},
			want: map[string]interface{}{
				"int32ToStr": map[string]interface{}{
					"-101": "-101",
					"0":    "zero",
					"255":  "0xff",
				},
			},
		}, {
			desc: "map fields 3",
			input: &pb3.Maps{
				Uint64ToEnum: map[uint64]pb3.Enum{
					1:  pb3.Enum_ONE,
					2:  pb3.Enum_TWO,
					10: pb3.Enum_TEN,
					47: 47,
				},
			},
			want: map[string]interface{}{
				"uint64ToEnum": map[string]interface{}{
					"1":  "ONE",
					"2":  "TWO",
					"10": "TEN",
					"47": int64(47),
				},
			},
		}, {
			desc: "map fields 4",
			input: &pb3.Maps{
				StrToNested: map[string]*pb3.Nested{
					"nested": {
						SString: "nested in a map",
					},
				},
			},
			want: map[string]interface{}{
				"strToNested": map[string]interface{}{
					"nested": map[string]interface{}{
						"sString": "nested in a map",
					},
				},
			},
		}, {
			desc: "map fields 5",
			input: &pb3.Maps{
				StrToOneofs: map[string]*pb3.Oneofs{
					"string": {
						Union: &pb3.Oneofs_OneofString{
							OneofString: "hello",
						},
					},
					"nested": {
						Union: &pb3.Oneofs_OneofNested{
							OneofNested: &pb3.Nested{
								SString: "nested oneof in map field value",
							},
						},
					},
				},
			},
			want: map[string]interface{}{
				"strToOneofs": map[string]interface{}{
					"nested": map[string]interface{}{
						"oneofNested": map[string]interface{}{
							"sString": "nested oneof in map field value",
						},
					},
					"string": map[string]interface{}{
						"oneofString": "hello",
					},
				},
			},
		}, {
			desc: "map field contains nil value",
			input: &pb3.Maps{
				StrToNested: map[string]*pb3.Nested{
					"nil": nil,
				},
			},
			want: map[string]interface{}{},
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
