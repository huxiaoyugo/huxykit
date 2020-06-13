package tlog

import (
	"bytes"
	"fmt"
	"time"
)

type logMsg struct {
	level     LogLevel
	msg       []interface{}
}

func (l *logMsg) ToBytes() []byte {
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("[%s] %s %v\n", l.level.String(),time.Now().Format("2006-01-02 15:04:05"), fmt.Sprint(l.msg...)))
	return b.Bytes()
}

type LogItem interface {
	ToBytes() []byte
}
