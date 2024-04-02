package zap

import (
	"fmt"
	"net/url"

	zaplib "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

func NewGrayLogCore(graylogURL string, facility string, level Level) (zapcore.Core, error) {
	graylogURI, err := url.Parse(graylogURL)
	if err != nil {
		return nil, err
	}

	var writerSyncer zapcore.WriteSyncer

	switch graylogURI.Scheme {
	case "udp":
		udpWriter, err := gelf.NewUDPWriter(graylogURI.Host)
		if err != nil {
			return nil, err
		}
		udpWriter.Facility = facility
		writerSyncer = zapcore.AddSync(udpWriter)
	case "tcp":
		tcpWriter, err := gelf.NewTCPWriter(graylogURI.Host)
		if err != nil {
			return nil, err
		}
		tcpWriter.Facility = facility
		writerSyncer = zapcore.AddSync(tcpWriter)
	}

	encoder := zapcore.NewJSONEncoder(zaplib.NewProductionEncoderConfig())

	levelEnablerFunc := zaplib.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return zapcore.Level(level) <= lvl
	})

	core := zapcore.NewCore(encoder, writerSyncer, levelEnablerFunc)
	err = core.Sync()
	if err != nil {
		return nil, fmt.Errorf("got error (%v) on creating graylog core", err)
	}

	return core, nil
}
