package jsonpb

import (
	"bytes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"io"
)

type JsonMarshal interface {
	JsonMarshal(writer io.Writer, opts Options)
}

type Student struct {
	Name        string
	Age         int
	Status      *wrappers.Int32Value  `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Msg         *wrappers.StringValue `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Class       *Class
	Details     []*Class
	ParamsOneof isParams_ParamsOneof
}

type isParams_ParamsOneof interface {
	isParams_ParamsOneof()
}

func (OneOfStruct) isParams_ParamsOneof() {}

type OneOfStruct struct {
	OneName   string
	OneAge    int
	OneStatus *wrappers.Int32Value  `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	OneMsg    *wrappers.StringValue `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (s *OneOfStruct) JsonMarshal(writer io.Writer, opts Options) {
	buf := &bytes.Buffer{}
	buf.WriteByte('{')
	WriteStructField(buf, NewString("one_name", s.OneName, opts), false)
	WriteStructField(buf, NewInt("one_age", int64(s.OneAge), opts), false)
	WriteStructField(buf, NewWrapperInt32("one_status", s.OneStatus, opts), false)
	WriteStructField(buf, NewWrapperString("one_msg", s.OneMsg, opts), true)
	buf.WriteByte('}')
	_, e := writer.Write(buf.Bytes())
	if e != nil {
		panic(e)
	}
}

func (s *Student) JsonMarshal(writer io.Writer, opts Options) {
	buf := &bytes.Buffer{}
	WriteObjectStart(buf)
	WriteStructField(buf, NewString("name", s.Name, opts), false)
	WriteStructField(buf, NewInt("age", int64(s.Age), opts), false)
	WriteStructField(buf, NewInterface("class", s.Class, opts), false)
	WriteStructField(buf, NewWrapperInt32("status", s.Status, opts), false)
	WriteStructField(buf, NewWrapperString("msg", s.Msg, opts), false)

	buf.WriteString(`"details":[`)
	for i := 0;i<len(s.Details);i++ {
		if i == len(s.Details) - 1 {
			WriteStructField(buf, NewInterface("", s.Details[i], opts), true)
		} else {
			WriteStructField(buf, NewInterface("", s.Details[i], opts), false)
		}
	}


	buf.WriteByte(']')
	buf.WriteByte(',')

	WriteStructField(buf, NewOneOf("params", s.ParamsOneof, opts), true)
	WriteObjectEnd(buf)
	WriteToWriter(writer, buf)
}

type Class struct {
	Name string
	No   int
}

func (s *Class) JsonMarshal(writer io.Writer, opts Options) {
	buf := &bytes.Buffer{}
	WriteObjectStart(buf)
	WriteStructField(buf, NewString("name", s.Name, opts), false)
	WriteStructField(buf, NewInt("no", int64(s.No), opts), true)
	WriteObjectEnd(buf)
	WriteToWriter(writer, buf)
}

func WriteStructField(buf *bytes.Buffer, typ FieldType, lastField bool) {
	typ.ToString(buf, lastField)
	if lastField {
		b := buf.Bytes()
		if len(b) > 0 {
			if b[len(b)-1] == ',' {
				b = b[:len(b)-1]
				buf.Reset()
				buf.Write(b)
			}
		}
	}
}

func WriteObjectStart(buf *bytes.Buffer) {
	buf.WriteByte('{')
}

func WriteObjectEnd(buf *bytes.Buffer) {
	buf.WriteByte('}')
}

func WriteToWriter(writer io.Writer, buf *bytes.Buffer) {
	_, e := writer.Write(buf.Bytes())
	if e != nil {
		panic(e)
	}
}