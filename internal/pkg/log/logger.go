package log

import (
	"errors"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"

	logFormatText = "text"
	logFormatJSON = "json"
)

var (
	_ Logger = logger{}
	_ Logger = entry{}

	Log logger
)

type (
	// Logger is what any Tendermint library should take.
	Logger interface {
		Debug(msg string, keyvals ...interface{})
		Info(msg string, keyvals ...interface{})
		Error(msg string, keyvals ...interface{})

		With(keyvals ...interface{}) Logger
	}

	logger struct {
		*logrus.Logger
	}
)

func init() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	Log = logger{log}
}

//SetLogger set a instance of Logger
func SetLogger(format, level string) {
	Log.SetFormatter(parseFormatter(format))
	Log.SetLevel(parseLevel(level))
}

func (l logger) Debug(msg string, keyvals ...interface{}) {
	l.WithFields(argsToFields(keyvals...)).Debug(msg)
}

func (l logger) Info(msg string, keyvals ...interface{}) {
	l.WithFields(argsToFields(keyvals...)).Info(msg)
}

func (l logger) Error(msg string, keyvals ...interface{}) {
	l.WithFields(argsToFields(keyvals...)).Error(msg)
}

func (l logger) With(keyvals ...interface{}) Logger {
	return entry{
		l.WithFields(argsToFields(keyvals...)),
	}
}

func (l logger) Write(p []byte) (n int, err error) {
	if len(p) > 1024 {
		l.Debug(string(p[:1024]), "truncated", true)
		return 1024, errors.New("log line too long")
	}
	l.Debug(string(p))
	return len(p), nil
}

type entry struct {
	*logrus.Entry
}

func (e entry) Debug(msg string, keyvals ...interface{}) {
	e.Entry.WithFields(argsToFields(keyvals...)).Debug(msg)
}

func (e entry) Info(msg string, keyvals ...interface{}) {
	e.Entry.WithFields(argsToFields(keyvals...)).Info(msg)
}

func (e entry) Error(msg string, keyvals ...interface{}) {
	e.Entry.WithFields(argsToFields(keyvals...)).Error(msg)
}

func (e entry) With(keyvals ...interface{}) Logger {
	return entry{
		e.WithFields(argsToFields(keyvals...)),
	}
}

func (e entry) Write(p []byte) (n int, err error) {
	return e.Write(p)
}

func argsToFields(keyvals ...interface{}) logrus.Fields {
	var fields = make(logrus.Fields)
	if len(keyvals)%2 != 0 {
		return fields
	}

	for i := 0; i < len(keyvals); i += 2 {
		fields[keyvals[i].(string)] = keyvals[i+1]
	}
	return fields
}

func parseFormatter(format string) logrus.Formatter {
	var formatter logrus.Formatter
	switch strings.ToLower(format) {
	case logFormatText:
		formatter = &logrus.TextFormatter{FullTimestamp: true}
	case logFormatJSON:
		formatter = &logrus.JSONFormatter{}
	default:
		formatter = &logrus.TextFormatter{FullTimestamp: true}
	}
	return formatter
}

func parseLevel(lvl string) logrus.Level {
	var level logrus.Level
	switch strings.ToLower(lvl) {
	case DebugLevel:
		level = logrus.DebugLevel
	case InfoLevel:
		level = logrus.InfoLevel
	case WarnLevel:
		level = logrus.WarnLevel
	case ErrorLevel:
		level = logrus.ErrorLevel
	default:
		level = logrus.InfoLevel
	}
	return level
}
