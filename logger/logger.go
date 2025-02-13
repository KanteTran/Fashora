package logger

import (
	"sync"

	"go.uber.org/zap"
)

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]any

var (
	once sync.Once
)

const (
	// Debug has verbose message
	debugLvl = "debug"
	//Info is default log level
	infoLvl = "info"
	// Warn is for logging messages about possible issues
	warnLvl = "warn"
	// Error is for logging errors
	errorLvl = "error"
	// Fatal is for logging fatal messages. The sytem shutsdown after logging the message.
	fatalLvl = "fatal"
)

// Logger is our contract for the logger
type Logger interface {
	Debug(msg string)
	Debugf(format string, args ...any)

	Info(msg string)
	Infof(format string, args ...any)
	InfoT(format string, args ...any)

	Warn(msg string)
	Warnf(format string, args ...interface{})

	Error(msg string)
	Errorf(format string, args ...any)
	ErrorT(format string, args ...any)

	Fatal(msg string)
	Fatalf(format string, args ...any)

	Panic(msg string)
	Panicf(format string, args ...any)

	WithFields(keyValues Fields) Logger

	GetDelegate() any
}

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

var log Logger = DefaultLogger()

func DefaultLogger() Logger {

	cfg := LoggerDefaultConfig()

	logger, _ := newZapLogger(cfg)
	return logger
}

// InitLogger returns an instance of logger
func InitLogger(config LoggerConfig) (Logger, error) {
	var err error
	once.Do(func() {
		log, err = newZapLogger(config)
	})
	return log, err
}

func Debug(msg string) {
	log.Debugf(msg)
}

func Debugf(format string, args ...any) {
	log.Debugf(format, args...)
}

func Info(msg string) {
	log.Infof(msg)
}

func Infof(format string, args ...any) {
	log.Infof(format, args...)
}

// InfoT ... stands for Info Terminate, it same as Infof()
// but we use it when logic flow is going to terminate after logging
func InfoT(format string, args ...any) {
	log.Infof("-----> "+format+"\n", args...)
}

func Warn(msg string) {
	log.Warnf(msg)
}

func Warnf(format string, args ...any) {
	log.Warnf(format, args...)
}

func Error(msg string) {
	log.Errorf(msg)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// ErrorT ... stands for Error Terminate, it same as Errorf()
// but we use it when logic flow is going to terminate after logging
func ErrorT(format string, args ...interface{}) {
	log.Errorf("-----> "+format+"\n", args...)
}

func Fatal(msg string) {
	log.Fatalf(msg)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panic(msg string) {
	log.Panicf(msg)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}

func GetDelegate() any {
	return log.GetDelegate()
}

func GetLogger() Logger {
	return log
}
