package log

import (
	"bytes"
	"fmt"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func Infof(format string, args ...any) {
	log.Infof(format, args...)
}

func Debugf(format string, args ...any) {
	log.Debugf(format, args...)
}

func Warnf(format string, args ...any) {
	log.Warnf(format, args...)
}

func Panicf(format string, args ...any) {
	log.Panicf(format, args...)
}

func Info(args ...any) {
	log.Info(args...)
}

func Fatal(args ...any) {
	log.Fatal(args...)
}

func Panic(args ...any) {
	log.Panic(args...)
}
func Debug(args ...any) {
	log.Debug(args...)
}

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&MyFormatter{})
}

type MyFormatter struct{}

func (m *MyFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[%s] [%s] [%s:%d %s] %s\n",
			timestamp, entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}
