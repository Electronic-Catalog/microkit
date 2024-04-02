package zap

import (
	"github.com/Electronic-Catalog/microkit/logger"
	"github.com/Electronic-Catalog/microkit/logger/keyval"
	zaplib "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"unsafe"
)

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

var DefaultStdLogger logger.Logger
var NopLogger logger.Logger

func init() {
	stdCore, _ := NewStandardCore(false, DebugLevel)
	zLogger := zaplib.New(stdCore)
	DefaultStdLogger = NewZapLogger(zLogger)
	NopLogger = NewZapLoggerWithCores(zapcore.NewNopCore())
}

type zapLogger struct {
	logger *zaplib.Logger
}

func NewZapLogger(logger *zaplib.Logger) logger.Logger {
	l := zapLogger{
		logger: logger,
	}

	return &l
}

func NewZapLoggerWithCores(cores ...zapcore.Core) logger.Logger {
	mergedCore := zapcore.NewTee(cores...)
	l := zapLogger{
		logger: zaplib.New(mergedCore),
	}

	return &l
}

func (l *zapLogger) Debug(message string, keyAndValues ...keyval.Pair) {
	var zapFields []zapcore.Field = *(*[]zapcore.Field)(unsafe.Pointer(&keyAndValues))
	l.logger.Debug(message, zapFields...)
}

func (l *zapLogger) Info(message string, keyAndValues ...keyval.Pair) {
	var zapFields []zapcore.Field = *(*[]zapcore.Field)(unsafe.Pointer(&keyAndValues))
	l.logger.Info(message, zapFields...)
}

func (l *zapLogger) Warn(message string, keyAndValues ...keyval.Pair) {
	var zapFields []zapcore.Field = *(*[]zapcore.Field)(unsafe.Pointer(&keyAndValues))
	l.logger.Warn(message, zapFields...)
}

func (l *zapLogger) Error(message string, keyAndValues ...keyval.Pair) {
	var zapFields []zapcore.Field = *(*[]zapcore.Field)(unsafe.Pointer(&keyAndValues))
	l.logger.Error(message, zapFields...)
}

func (l *zapLogger) Panic(message string, keyAndValues ...keyval.Pair) {
	var zapFields []zapcore.Field = *(*[]zapcore.Field)(unsafe.Pointer(&keyAndValues))
	l.logger.Panic(message, zapFields...)
}
