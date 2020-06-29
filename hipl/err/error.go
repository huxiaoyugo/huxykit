package err

import (
	"bytes"
	"fmt"
)

type Error interface {
	error
	Append(...string)Error
}

type IError struct {
	errInfo []string
}

func (p *IError) Append(msg ...string) Error{
	p.errInfo = append(p.errInfo, msg...)
	return p
}

func (p *IError) Error() string {
	buf := bytes.NewBufferString("")
	for _, msg := range p.errInfo {
		buf.WriteString(msg)
		buf.WriteByte(' ')
	}
	return buf.String()
}

var _ Error = &IError{}



func Err(msg ...string) Error {
	return (&IError{}).Append(msg...)
}

func ErrF(str string, args ... interface{}) Error {
	return Err(fmt.Sprintf(str, args...))
}