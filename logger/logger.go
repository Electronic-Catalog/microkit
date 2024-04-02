package logger

import "github.com/Electronic-Catalog/microkit/logger/keyval"

type Logger interface {
	Debug(message string, keyAndValues ...keyval.Pair)
	Info(message string, keyAndValues ...keyval.Pair)
	Warn(message string, keyAndValues ...keyval.Pair)
	Error(message string, keyAndValues ...keyval.Pair)
	Panic(message string, keyAndValues ...keyval.Pair)
}
