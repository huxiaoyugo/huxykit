package tlog

import (
	"bytes"
	"fmt"
	"github.com/huxiaoyugo/huxykit/utils"
	"io"
	"os"
	"strings"
	"time"
)

const (
	// 默认切分周期
	defaultRotateCycle = CycleMin
	// 默认日志路径
	defaultLogPath = "./log"
	// 默认的文件名
	defaultLogName = "app.log"
	// wf日志文件名后缀
	suffix = ".wf"
)

const (
	CycleHalfHour Cycle = iota
	CycleHour
	CycleDay
	CycleMin
)

const (
	logTypeNormal logType = iota
	logTypeWf
)

// 周期类型
type Cycle int

type logType int

// 周期的持续时间
func (c Cycle) DurationSec() int64 {
	switch c {
	case CycleHalfHour:
		return 1800
	case CycleHour:
		return 3600
	case CycleDay:
		return 86400
	case CycleMin:
		return 60
	default:
		// 默认返回一个小时
		return 3600
	}
}

func (c Cycle) RotateIndex() string {
	switch c {
	case CycleHour, CycleHalfHour:
		interval := utils.GetIntervalFromDawn(time.Now())
		index := int64(interval.Seconds()) / c.DurationSec()
		return fmt.Sprintf("%s%02d", time.Now().Format("20060102"), index)
	case CycleMin:
		interval := utils.GetIntervalFromDawn(time.Now())
		index := int64(interval.Seconds()) / c.DurationSec()
		return fmt.Sprintf("%s%04d", time.Now().Format("20060102"), index)
	case CycleDay:
		return time.Now().Format("20060102")
	}
	return ""
}

type FileWriter struct {
	// 是否需要确切分日志文件
	rotateSwitch bool
	// 切分周期
	rotateCycle Cycle
	// 文件路径
	logPath string
	// 文件名
	logName string
	// 上次切分文件的时间戳
	lastRotateTime int64
	// 文件切分索引如2020062801
	curRotateIndex string

	// 正常的日志
	file *os.File
	// wf日志
	fileWf *os.File
}

type FileWriterConfig struct {
	RotateSwitch bool
	RotateCycle  Cycle
	Path         string
	LogName      string
}

func DefaultFileWriter() *FileWriter {
	config := FileWriterConfig{
		RotateSwitch: true,
		RotateCycle:  defaultRotateCycle,
		Path:         defaultLogPath,
		LogName:      defaultLogName,
	}
	return NewFileWriter(config)
}

func NewFileWriter(config FileWriterConfig) *FileWriter {
	w := &FileWriter{
		rotateSwitch: config.RotateSwitch,
		rotateCycle:  config.RotateCycle,
		logPath:      config.Path,
		logName:      config.LogName,
	}
	// 创建logPath
	_ = w.createLogPath()
	return w
}

func (f *FileWriter) Write(p []byte) (n int, err error) {

	if err = f.setLogFile(); err != nil {
		return
	}
	switch f.logType(p) {
	case logTypeWf:
		return f.fileWf.Write(p)
	default:
		return f.file.Write(p)
	}
}

func (f *FileWriter) logType(p []byte) logType {
	t := logTypeNormal
	if bytes.HasPrefix(p, []byte("[ERROR]")) ||
		bytes.HasPrefix(p, []byte("[WARN]")) ||
		bytes.HasPrefix(p, []byte("[FATAL]")) {
		t = logTypeWf
	}
	return t
}

func (f *FileWriter) setLogFile() error {
	var err error
	ri := strings.Trim(f.rotateCycle.RotateIndex(), " ")
	if ri == "" {
		return fmt.Errorf("rotateIndex can not be nil")
	}
	if ri == f.curRotateIndex {
		return nil
	}

	// 关闭原有的os.File
	f.closeFile()

	// 修改app.log和app.log.wf文件名
	err = f.rename(f.curRotateIndex)
	if err != nil {
		return err
	}

	// 修改当前的index
	f.curRotateIndex = ri

	// 重新打开
	if f.file, err = f.openFile(f.logName); err != nil {
		return err
	}
	if f.fileWf, err = f.openFile(f.logName + suffix); err != nil {
		return err
	}
	return nil
}

// 打开文件
func (f *FileWriter) openFile(fileName string) (*os.File, error) {
	file := f.logPath
	if !strings.HasSuffix(file, "/") {
		file += "/"
	}
	file += fileName
	fn, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fn, err
	}
	err = fn.Chmod(0666)
	if err != nil {
		return nil, err
	}
	return fn, nil
}

// 关闭file和fileWf
func (f *FileWriter) closeFile() {
	if f.fileWf != nil {
		f.fileWf.Close()
		f.fileWf = nil
	}
	if f.file != nil {
		f.file.Close()
		f.file = nil
	}
}

// 当需要切分日志的时候
// 将app.log ==> app.log.2020010201
// 将app.log.wf ==> app.log.wf.2020010201
func (f *FileWriter) rename(rotateIndex string) error {
	rotateIndex = strings.Trim(rotateIndex, " ")
	if rotateIndex == "" {
		return nil
	}

	// rename app.log
	oldFile := f.logPath
	if !strings.HasSuffix(oldFile, "/") {
		oldFile += "/"
	}
	oldFile += f.logName
	newFile := fmt.Sprintf("%s.%s", oldFile, rotateIndex)
	if utils.FileExist(oldFile) {
		if err := os.Rename(oldFile, newFile); err != nil {
			return err
		}
	}

	// rename app.log.wf => app.log.wf.2020060201
	oldFile += suffix
	newFile = fmt.Sprintf("%s.%s", oldFile, rotateIndex)
	if utils.FileExist(oldFile) {
		if err := os.Rename(oldFile, newFile); err != nil {
			return err
		}
	}
	return nil
}

// 创建logPath的文件夹
func (f *FileWriter) createLogPath() error {
	if !utils.FileExist(f.logPath) {
		err := os.Mkdir(f.logPath, 0777)
		if err != nil {
			return err
		}
	}
	return os.Chmod(f.logPath, 0777)
}

var _ io.Writer = &FileWriter{}
