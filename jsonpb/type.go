package jsonpb

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"io"
	"reflect"
)

type FieldType interface {
	ToString(writer io.Writer, lastField bool)
	GetTypeName() string
}

type BaseFieldType struct {
	Name     string
	Opts     Options
}

func (b BaseFieldType) ToString(writer io.Writer, lastField bool) {
	if b.Name == "" {
		return
	}
	_, e := writer.Write([]byte{'"'})
	if e != nil {
		panic(e)
	}
	_, e = writer.Write([]byte(b.Name))
	if e != nil {
		panic(e)
	}
	_, e = writer.Write([]byte("\":"))
	if e != nil {
		panic(e)
	}

}

func (b BaseFieldType) GetTypeName() string {
	return b.Name
}

var _ FieldType = BaseFieldType{}

// ================= String ==============
type String struct {
	BaseFieldType
	Val string
}

func (b String) ToString(writer io.Writer, lastField bool) {
	if b.Opts.IgnoreEmpty && b.Val == "" {
		return
	}
	b.BaseFieldType.ToString(writer,lastField)
	_, e := writer.Write([]byte{'"'})
	if e != nil {
		panic(e)
	}
	_, e = writer.Write([]byte(b.Val))
	if e != nil {
		panic(e)
	}
	_, e = writer.Write([]byte{'"'})
	if e != nil {
		panic(e)
	}
	if !lastField {
		_, e = writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewString(name string, val string, opts Options) String {
	return String{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// =================  Int ==============
type Int struct {
	BaseFieldType
	Val int64
}

func (b Int) ToString(writer io.Writer, lastField bool) {
	if b.Opts.IgnoreEmpty && b.Val == 0 {
		return
	}
	b.BaseFieldType.ToString(writer, lastField)
	// todo: 优化
	v := fmt.Sprint(b.Val)
	_, e := writer.Write([]byte(v))
	if e != nil {
		panic(e)
	}
	if !lastField {
		_, e = writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewInt(name string, val int64, opts Options) Int {
	return Int{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// =================  UInt ==============
type UInt struct {
	BaseFieldType
	Val int64
}

func (b UInt) ToString(writer io.Writer, lastField bool) {
	if b.Opts.IgnoreEmpty && b.Val == 0 {
		return
	}
	b.BaseFieldType.ToString(writer,lastField)
	// todo: 优化
	v := fmt.Sprint(b.Val)
	_, e := writer.Write([]byte(v))
	if e != nil {
		panic(e)
	}
	if !lastField {
		_, e = writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewUInt(name string, val int64, opts Options) UInt {
	return UInt{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// =================  Interface ==============
type Interface struct {
	BaseFieldType
	Val interface{}
}

func (b Interface) ToString(writer io.Writer, lastField bool) {

	// todo: 处理0值
	if reflect.TypeOf(b).Kind() == reflect.Ptr {
		if reflect.ValueOf(b.Val).IsNil() {
			return
		}
	}

	if b.Opts.IgnoreEmpty && b.Val == "" {
		return
	}
	b.BaseFieldType.ToString(writer, lastField)
	if jm, ok := b.Val.(JsonMarshal); ok {
		jm.JsonMarshal(writer, b.Opts)
		if !lastField {
			_, e := writer.Write([]byte{','})
			if e != nil {
				panic(e)
			}
		}
		return
	}
	panic(fmt.Errorf("%s do not implement interface 'JsonMarshal'", reflect.TypeOf(b.Val).Name()))
}

func NewInterface(name string, val interface{}, opts Options) Interface {
	return Interface{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// =================  Float ==============
type Float struct {
	BaseFieldType
	Val float64
}

func (b Float) ToString(writer io.Writer, lastField bool) {
	if b.Opts.IgnoreEmpty && b.Val == 0 {
		return
	}
	b.BaseFieldType.ToString(writer, lastField)
	// todo: 优化
	v := fmt.Sprint(b.Val)
	_, e := writer.Write([]byte(v))
	if e != nil {
		panic(e)
	}
	if !lastField {
		_, e = writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewFloat(name string, val float64, opts Options) Float {
	return Float{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// ================= OneOf ==============
type OneOf struct {
	BaseFieldType
	Val interface{}
}

func (b OneOf) ToString(writer io.Writer, lastField bool) {
	buf := &bytes.Buffer{}
	NewInterface("", b.Val, b.Opts).ToString(buf, true)
	data := buf.Bytes()
	// 去掉{}
	if len(data) >= 2 {
		data = data[1 : len(data)-1]
	}
	_, e := writer.Write(data)
	if e != nil {
		panic(e)
	}
	if !lastField {
		_, e = writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewOneOf(name string, val interface{}, opts Options) OneOf {
	return OneOf{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// ================= WrapperString ==============
type WrapperString struct {
	BaseFieldType
	Val *wrappers.StringValue
}

func (b WrapperString) ToString(writer io.Writer, lastField bool) {
	if b.Opts.IgnoreEmpty && b.Val.GetValue() == "" {
		return
	}
	b.BaseFieldType.ToString(writer, lastField)
	_, e := writer.Write([]byte{'"'})
	if e != nil {
		panic(e)
	}
	_, e = writer.Write([]byte(b.Val.GetValue()))
	if e != nil {
		panic(e)
	}
	_, e = writer.Write([]byte{'"'})
	if e != nil {
		panic(e)
	}
	if !lastField {
		_, e = writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewWrapperString(name string, val *wrappers.StringValue, opts Options) WrapperString {
	return WrapperString{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// =================  WrapperBool ==============
type WrapperBool struct {
	BaseFieldType
	Val *wrappers.BoolValue
}

func (b WrapperBool) ToString(writer io.Writer, lastField bool) {
	if b.Opts.IgnoreEmpty && b.Val == nil {
		return
	}
	b.BaseFieldType.ToString(writer, lastField)
	if b.Val.GetValue() {
		_, e := writer.Write([]byte("true"))
		if e != nil {
			panic(e)
		}
	} else {
		_, e := writer.Write([]byte("false"))
		if e != nil {
			panic(e)
		}
	}
	if !lastField {
		_, e := writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewWrapperBool(name string, val *wrappers.BoolValue, opts Options) WrapperBool {
	return WrapperBool{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// =================  WrapperInt32 ==============
type WrapperInt32 struct {
	BaseFieldType
	Val *wrappers.Int32Value
}

func (b WrapperInt32) ToString(writer io.Writer, lastField bool) {
	if b.Opts.IgnoreEmpty && b.Val.GetValue() == 0 {
		return
	}
	b.BaseFieldType.ToString(writer, lastField)
	// todo
	v := fmt.Sprint(b.Val.GetValue())
	_, e := writer.Write([]byte(v))
	if e != nil {
		panic(e)
	}
	if !lastField {
		_, e = writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewWrapperInt32(name string, val *wrappers.Int32Value, opts Options) WrapperInt32 {
	return WrapperInt32{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// =================  WrapperUInt32 ==============
type WrapperUInt32 struct {
	BaseFieldType
	Val *wrappers.UInt32Value
}

func (b WrapperUInt32) ToString(writer io.Writer, lastField bool) {
	if b.Opts.IgnoreEmpty && b.Val.GetValue() == 0 {
		return
	}
	b.BaseFieldType.ToString(writer, lastField)
	// todo
	v := fmt.Sprint(b.Val.GetValue())
	_, e := writer.Write([]byte(v))
	if e != nil {
		panic(e)
	}
	if !lastField {
		_, e = writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewWrapperUInt32(name string, val *wrappers.UInt32Value, opts Options) WrapperUInt32 {
	return WrapperUInt32{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// =================  WrapperInt64 ==============
type WrapperInt64 struct {
	BaseFieldType
	Val *wrappers.Int64Value
}

func (b WrapperInt64) ToString(writer io.Writer, lastField bool) {
	if b.Opts.IgnoreEmpty && b.Val.GetValue() == 0 {
		return
	}
	b.BaseFieldType.ToString(writer, lastField)
	// todo
	v := fmt.Sprint(b.Val.GetValue())
	_, e := writer.Write([]byte(v))
	if e != nil {
		panic(e)
	}
	if !lastField {
		_, e = writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewWrapperInt64(name string, val *wrappers.Int64Value, opts Options) WrapperInt64 {
	return WrapperInt64{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// =================  WrapperUInt64 ==============
type WrapperUInt64 struct {
	BaseFieldType
	Val *wrappers.UInt64Value
}

func (b WrapperUInt64) ToString(writer io.Writer, lastField bool) {
	if b.Opts.IgnoreEmpty && b.Val.GetValue() == 0 {
		return
	}
	b.BaseFieldType.ToString(writer, lastField)
	// todo
	v := fmt.Sprint(b.Val.GetValue())
	_, e := writer.Write([]byte(v))
	if e != nil {
		panic(e)
	}
	if !lastField {
		_, e = writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewWrapperUInt64(name string, val *wrappers.UInt64Value, opts Options) WrapperUInt64 {
	return WrapperUInt64{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// =================  WrapperFloat ==============
type WrapperFloat struct {
	BaseFieldType
	Val *wrappers.FloatValue
}

func (b WrapperFloat) ToString(writer io.Writer, lastField bool) {
	if b.Opts.IgnoreEmpty && b.Val.GetValue() == 0 {
		return
	}
	b.BaseFieldType.ToString(writer, lastField)
	// todo
	v := fmt.Sprint(b.Val.GetValue())
	_, e := writer.Write([]byte(v))
	if e != nil {
		panic(e)
	}
	if !lastField {
		_, e = writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewWrapperFloat(name string, val *wrappers.FloatValue, opts Options) WrapperFloat {
	return WrapperFloat{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// =================  WrapperDouble ==============
type WrapperDouble struct {
	BaseFieldType
	Val *wrappers.DoubleValue
}

func (b WrapperDouble) ToString(writer io.Writer, lastField bool) {
	if b.Opts.IgnoreEmpty && b.Val.GetValue() == 0 {
		return
	}
	b.BaseFieldType.ToString(writer, lastField)
	// todo
	v := fmt.Sprint(b.Val.GetValue())
	_, e := writer.Write([]byte(v))
	if e != nil {
		panic(e)
	}
	if !lastField {
		_, e = writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewWrapperDouble(name string, val *wrappers.DoubleValue, opts Options) WrapperDouble {
	return WrapperDouble{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}

// =================  WrapperBytes ==============
type WrapperBytes struct {
	BaseFieldType
	Val *wrappers.BytesValue
}

func (b WrapperBytes) ToString(writer io.Writer, lastField bool) {
	if b.Opts.IgnoreEmpty && len(b.Val.GetValue()) == 0 {
		return
	}
	// base64
	b.BaseFieldType.ToString(writer, lastField)

	_, e := writer.Write([]byte{'"'})
	if e != nil {
		panic(e)
	}
	_, e = writer.Write([]byte(base64.StdEncoding.EncodeToString(b.Val.GetValue())))
	if e != nil {
		panic(e)
	}
	_, e = writer.Write([]byte{'"'})
	if e != nil {
		panic(e)
	}
	if !lastField {
		_, e = writer.Write([]byte{','})
		if e != nil {
			panic(e)
		}
	}
}

func NewWrapperBytes(name string, val *wrappers.BytesValue, opts Options) WrapperBytes {
	return WrapperBytes{
		BaseFieldType: BaseFieldType{
			Name:     name,
			Opts:     opts,
		},
		Val: val,
	}
}
