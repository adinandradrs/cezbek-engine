package main

import (
	"github.com/adinandradrs/cezbek-engine/internal"
	"go.uber.org/zap"
)

func main() {
	c := internal.NewContainer("APP_CEZBEK_API")
	p := c.LoadPool()
	c.Logger.Info("pool database", zap.Any("", p.Pool))
}
