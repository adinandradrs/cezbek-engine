package apps

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLog(prod bool) (z *zap.Logger, fh fiber.Handler) {
	var c zap.Config
	defCfg := logger.ConfigDefault
	defCfg.Format = `${time} {"router_activity" : [${status},"${latency}","${method}","${path}"]}` + "\n"
	defCfg.TimeFormat = "2006-01-02T15:04:05.000Z0700"
	httpLogger := logger.New(defCfg)
	if prod {
		c = zap.NewProductionConfig()
		c.DisableStacktrace = true

	} else {
		c = zap.NewDevelopmentConfig()
		c.DisableStacktrace = false
	}
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	z, _ = c.Build()
	return z, httpLogger
}
