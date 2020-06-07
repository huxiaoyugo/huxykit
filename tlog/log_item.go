package tlog

import (
	"bytes"
	"fmt"
)

type logMsg struct {
	level LogLevel
	msg  []interface{}
}

func(l *logMsg) ToBytes() []byte {
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("[%s] %s", l.level.String(), fmt.Sprint(l.msg...)))
	return b.Bytes()
}

type LogItem interface {
	ToBytes() []byte
}