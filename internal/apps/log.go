package apps

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLog(prod bool) (z *zap.Logger) {
	var c zap.Config
	if prod {
		c = zap.NewProductionConfig()
		c.DisableStacktrace = true

	} else {
		c = zap.NewDevelopmentConfig()
		c.DisableStacktrace = false
	}
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	z, _ = c.Build()
	return z

}
