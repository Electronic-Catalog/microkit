package zap

import (
	"fmt"

	sentry "github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap/zapcore"
)

func NewSentryCore(dsn string, tags map[string]string) (zapcore.Core, error) {
	zapConf := sentry.Configuration{
		Tags:  tags,
		Level: zapcore.ErrorLevel,
	}

	core, err := sentry.NewCore(zapConf, sentry.NewSentryClientFromDSN(dsn))
	if err != nil {
		return nil, fmt.Errorf("got error (%v) on creating sentry core", err)
	}

	return core, nil
}
