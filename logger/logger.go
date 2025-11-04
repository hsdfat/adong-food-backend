package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerI interface {

	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})

}


type Logger struct {
	*zap.SugaredLogger
}

func NewLogger() *Logger {
	// set caller skip to 2
	cfg := zap.NewProductionConfig()
	
	cfg.Level.SetLevel(zapcore.DebugLevel) // Set the desired level (e.g., InfoLevel)
	logger, _ := cfg.Build()
	logger = logger.WithOptions(zap.AddCallerSkip(1))
	
	sugar := logger.Sugar()
	
	return &Logger{
		SugaredLogger: sugar,
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.SugaredLogger.With(args...).Info(msg)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.SugaredLogger.With(args...).Warn(msg)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.SugaredLogger.With(args...).Error(msg)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.SugaredLogger.With(args...).Debug(msg)
}

var (
	Log  LoggerI = NewLogger()
)