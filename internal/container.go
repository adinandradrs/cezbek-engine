package internal

import (
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type (
	Container struct {
		Logger *zap.Logger
		Viper  *viper.Viper
		app    string
	}

	Env struct {
		ContextPath string
		HttpPort    string
	}
)

func NewContainer(app string) Container {
	logger := apps.NewLog(false)
	conf, err := apps.NewEnv(logger)
	if err != nil {
		logger.Panic("error to load config", zap.Any("", &err))
	}
	return svcRegister(&Container{
		Logger: logger,
		Viper:  conf,
		app:    app,
	})
}

func svcRegister(c *Container) Container {
	clientSvc := adaptor.NewConsul(adaptor.Consul{
		Host:    c.Viper.GetString("consul_host"),
		Port:    c.Viper.GetInt("consul_port"),
		Service: c.Viper.GetString(c.app),
		Viper:   c.Viper,
		Logger:  c.Logger,
	})
	ex := clientSvc.Register()
	if ex != nil {
		c.Logger.Panic("error to register", zap.Any("", &ex))
	}
	return *c
}

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
