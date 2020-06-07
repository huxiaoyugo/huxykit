package tlog

import (
	"fmt"
	"io"
)

const (
	logLevelDebug LogLevel = iota
	logLevelInfo
	logLevelWarn
	logLevelError
	logLevelFatal
)

type LogLevel int

func (l LogLevel) String() string {
	switch l {
	case logLevelDebug:
		return "DEBUG"
	case logLevelInfo:
		return "INFO"
	case logLevelWarn:
		return "WARN"
	case logLevelError:
		return "ERROR"
	case logLevelFatal:
		return "FATAL"
	}
	return ""
}

// FATAL、ERROR、WARN、INFO、DEBUG
type Logger interface {
	Debug(msg ...interface{})
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
	Fatal(msg ...interface{})
}

const (
	// 日志输出chan的buffer长度
	chanBufferLen = 200
	// 默认切分周期
	defaultRotateCycle = CycleHalfHour

	// 默认日志路径
	defaultLogPath = "./log"
)

const (
	CycleHalfHour Cycle = iota
	CycleHour
	CycleDay
	CycleWeek
)

// 周期类型
type Cycle int

type tLog struct {
	io.Writer
	// 日志输出级别
	level LogLevel
	// 是否需要确切分日志文件
	rotateSwitch bool
	// 切分周期
	rotateCycle Cycle
	// 文件路径
	logPath string
	//
	// 存放日志信息的chan
	buffer chan LogItem
}

func NewLogger(writer io.Writer, level LogLevel) Logger {
	logger := &tLog{
		Writer:       writer,
		level:        level,
		rotateSwitch: true,
		rotateCycle:  defaultRotateCycle,
		logPath:      defaultLogPath,
		buffer:       make(chan LogItem, chanBufferLen),
	}
	// 开始日志收集
	logger.Start()
	return logger
}

func(t *tLog) Start() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				 fmt.Println(err)
			}
		}()
		for {
			select {
			case logItem := <-t.buffer:
				_, _ = t.Writer.Write(logItem.ToBytes())
			}
		}
	}()
}

func (t *tLog) Debug(msg ...interface{}) {
	t.writeLog(logLevelDebug, msg...)
}

func (t *tLog) Info(msg ...interface{}) {
	t.writeLog(logLevelInfo, msg...)
}

func (t *tLog) Warn(msg ...interface{}) {
	t.writeLog(logLevelWarn, msg...)
}

func (t *tLog) Error(msg ...interface{}) {
	t.writeLog(logLevelError, msg...)
}

func (t *tLog) Fatal(msg ...interface{}) {
	t.writeLog(logLevelFatal, msg...)
}

func (t *tLog) writeLog(level LogLevel, msg ...interface{}) {
	logItem := &logMsg{
		level: level,
		msg:   msg,
	}
	select {
	case t.buffer <- logItem:
	default:
		// 防止日志输出阻塞主goroutine
		fmt.Println()
	}
}

var _ Logger = &tLog{}
