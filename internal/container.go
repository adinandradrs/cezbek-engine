package internal

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Container struct {
	Logger *zap.Logger
	Viper  *viper.Viper
}

func NewContainer() Container {
	log := apps.NewLog(false)
	vp, err := apps.NewEnv(log)
	if err != nil {
		log.Fatal("error to load config", zap.Any("", &err))
	}
	return Container{
		Logger: log,
		Viper:  vp,
	}
}

type (
	Env struct {
		ContextPath string
		HttpPort    string
	}
)

func (c *Container) LoadPool() *storage.PgPool {
	return storage.NewPgPool(&storage.PgOptions{
		Host:   c.Viper.GetString("db_host"),
		Port:   c.Viper.GetString("db_port"),
		User:   c.Viper.GetString("db_username"),
		Passwd: c.Viper.GetString("db_password"),
		Schema: c.Viper.GetString("db_schema"),
		Logger: c.Logger,
	})
}
