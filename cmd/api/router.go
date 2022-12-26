package main

import (
	"github.com/adinandradrs/cezbek-engine/internal"
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"strconv"
)

func main() {
	c := internal.NewContainer()
	c.Viper = registerService(&c)
	p := c.LoadPool()
	c.Logger.Info("pool database", zap.Any("", p.Pool))
}

func registerService(c *internal.Container) *viper.Viper {
	port, _ := strconv.Atoi(os.Getenv("CONSUL_PORT"))
	consulWatcher := adaptor.NewConsul(adaptor.Consul{
		Port:    port,
		Host:    os.Getenv("CONSUL_HOST"),
		Service: os.Getenv("APP_CEZBEK_API"),
		Viper:   c.Viper,
	})
	ex := consulWatcher.Register()
	if ex != nil {
		c.Logger.Panic("failed to register, stop immediate", zap.Any("", ex))
	}
	return c.Viper
}
