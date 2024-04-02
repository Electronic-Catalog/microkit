package zap

import (
	"github.com/Electronic-Catalog/microkit/logger/keyval"
	zaplib "go.uber.org/zap"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := DefaultStdLogger

	logger.Info("It works")
}

func BenchmarkLogger(t *testing.B) {
	logger := DefaultStdLogger
	t.StartTimer()
	for i := 0; i < t.N; i++ {
		logger.Info("test benchmark",
			keyval.String("fn", "test message"),
			keyval.Int("aNumber", 2000),
		)
	}
	t.StopTimer()
}

func BenchmarkZapLogger(t *testing.B) {
	zpCore, _ := NewStandardCore(false, InfoLevel)
	logger := zaplib.New(zpCore)

	t.StartTimer()
	for i := 0; i < t.N; i++ {
		logger.Info("benchmark zaplib", zaplib.String("fn", "test message"), zaplib.Int("aNumber", 2000))
	}
	t.StopTimer()
}
