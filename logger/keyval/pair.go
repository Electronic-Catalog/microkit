package keyval

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Pair zapcore.Field

func String(key string, val string) Pair {
	return Pair(zap.String(key, val))
}

func Int(key string, val int) Pair {
	return Pair(zap.Int(key, val))
}

func Int32(key string, val int32) Pair {
	return Pair(zap.Int32(key, val))
}

func Int64(key string, val int64) Pair {
	return Pair(zap.Int64(key, val))
}

func Float32(key string, val float32) Pair {
	return Pair(zap.Float32(key, val))
}

func Float64(key string, val float64) Pair {
	return Pair(zap.Float64(key, val))
}

func Binary(key string, val []byte) Pair {
	return Pair(zap.Binary(key, val))
}

func Error(err error) Pair {
	return Pair(zap.Error(err))
}
