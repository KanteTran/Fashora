package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getZapLevel(level string) zapcore.Level {
	switch level {
	case infoLvl:
		return zapcore.InfoLevel
	case warnLvl:
		return zapcore.WarnLevel
	case debugLvl:
		return zapcore.DebugLevel
	case errorLvl:
		return zapcore.ErrorLevel
	case fatalLvl:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l *zapLogger) Debug(msg string) {
	l.sugaredLogger.Debug(msg)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

// InfoT ... stands for Info Terminate, it same as Infof()
// but we use it when logic flow is going to terminate after logging
func (l *zapLogger) InfoT(format string, args ...interface{}) {
	l.sugaredLogger.Infof("-----> "+format+"\n", args...)
}

func (l *zapLogger) Info(msg string) {
	l.sugaredLogger.Info(msg)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

func (l *zapLogger) Warn(msg string) {
	l.sugaredLogger.Warn(msg)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

// ErrorT ... stands for Error Terminate, it same as Errorf()
// but we use it when logic flow is going to terminate after logging
func (l *zapLogger) ErrorT(format string, args ...interface{}) {
	l.sugaredLogger.Errorf("-----> "+format+"\n", args...)
}

func (l *zapLogger) Error(msg string) {
	l.sugaredLogger.Error(msg)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) Fatal(msg string) {
	l.sugaredLogger.Fatal(msg)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Panicf(format, args...)
}

func (l *zapLogger) Panic(msg string) {
	l.sugaredLogger.Panic(msg)

}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{newLogger}
}

func (l *zapLogger) GetDelegate() interface{} {
	return l.sugaredLogger
}

func getEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func newZapLogger(config LoggerConfig) (Logger, error) {
	cores := []zapcore.Core{}

	if config.EnableConsole {
		level := getZapLevel(config.ConsoleLevel)
		writer := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(getEncoder(config.ConsoleJSONFormat), writer, level)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	// AddCallerSkip skips 2 number of callers, this is important else the file that gets
	// logged will always be the wrapped file. In our case zap.go
	logger := zap.New(combinedCore,
		zap.AddCallerSkip(2),
		zap.AddCaller(),
	).Sugar()

	return &zapLogger{
		sugaredLogger: logger,
	}, nil
}
