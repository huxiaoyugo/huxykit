package jsonpb

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"os"
	"testing"
)

func TestJsonPb(t *testing.T) {

	s := Student{
		Name:  "huxiaoyu",
		Age:   20,
		Status: &wrappers.Int32Value{
			Value: 100,
		},
		Msg: &wrappers.StringValue{
			Value:"this is msg",
		},
		Class: &Class{
			Name: "math",
			No:   1,
		},
		Details: []*Class{
			{
				Name: "Chinese",
				No:   2,
			},
			{
				Name: "math",
				No:   1,
			},
		},
		ParamsOneof: &OneOfStruct{
			OneName:   "param_name",
			OneAge:    0,
			OneStatus: nil,
			OneMsg:    nil,
		},
	}
	s.JsonMarshal(os.Stdout, Options{
		PrettyMode:  "",
		Indent:      "",
		IgnoreEmpty: true,
	})
}
