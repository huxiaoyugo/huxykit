package tlog

import (
	"os"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	logger := NewDefaultLogger(logLevelDebug)

	for i:=0;;i++ {
		logger.Debug(i,"this is debug")
		logger.Error(i,"this is error")
		logger.Warn(i,"this is warn")
		logger.Fatal(i,"this is fatal")
		logger.Info(i,"this is info")
		time.Sleep(time.Second)
	}
}

func TestOpenFile(t *testing.T) {

	_, err := os.OpenFile("./log/apps.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(0666))
	if err != nil {
		t.Error(err)
		return
	}
}