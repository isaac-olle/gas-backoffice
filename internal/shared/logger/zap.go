package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLoggerAdapter struct {
	logger *zap.Logger
}

func NewZapLogger(env string) (*ZapLoggerAdapter, error) {
	var cfg zap.Config
	if env == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.DisableStacktrace = true
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLoggerAdapter{logger: logger}, nil
}

func (z *ZapLoggerAdapter) Info(msg string, fields ...Field) {
	z.logger.Info(msg, toZapFields(fields)...)
}

func (z *ZapLoggerAdapter) Error(msg string, fields ...Field) {
	z.logger.Error(msg, toZapFields(fields)...)
}

func (z *ZapLoggerAdapter) Fatal(msg string, fields ...Field) {
	z.logger.Fatal(msg, toZapFields(fields)...)
}

func (z *ZapLoggerAdapter) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, toZapFields(fields)...)
}

func toZapFields(fields []Field) []zap.Field {
	zf := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zf = append(zf, zap.Any(f.Key, f.Value))
	}
	return zf
}
