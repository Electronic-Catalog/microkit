package zap

import (
	"os"

	zaplib "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewStandardCore(pretty bool, level Level) (zapcore.Core, error) {
	writerSyncer := zapcore.AddSync(os.Stdout)

	var encoder zapcore.Encoder
	if pretty {
		encoder = zapcore.NewConsoleEncoder(zaplib.NewDevelopmentEncoderConfig())
	} else {
		encoder = zapcore.NewJSONEncoder(zaplib.NewProductionEncoderConfig())
	}

	levelEnablerFunc := zaplib.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return zapcore.Level(level) <= lvl
	})

	return zapcore.NewCore(encoder, writerSyncer, levelEnablerFunc), nil
}
