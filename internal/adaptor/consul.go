package adaptor

import (
	"fmt"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"strconv"
)

type Consul struct {
	Host    string
	Port    int
	Service string
	Logger  *zap.Logger
	Viper   *viper.Viper
}

type ConsulWatcher interface {
	Register() *model.TechnicalError
}

func NewConsul(c Consul) ConsulWatcher {
	return &c
}

func (c Consul) Register() *model.TechnicalError {
	client, err := consul.NewClient(&consul.Config{
		Address: c.Host + ":" + strconv.Itoa(c.Port),
	})
	if err != nil {
		c.Logger.Fatal("cannot initialize consul client", zap.Error(err))
	}
	err = client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      c.Service,
		Name:    c.Service,
		Address: c.Host,
		Port:    c.Port,
	})
	if err != nil {
		return apps.Exception("failed to register service", err, zap.Any("", c), c.Logger)
	}
	ex := c.readConfig()
	if ex != nil {
		return ex
	}
	return nil
}

func (c Consul) readConfig() *model.TechnicalError {
	err := c.Viper.AddRemoteProvider("consul",
		fmt.Sprintf("%s:%v", c.Host, c.Port),
		c.Service)
	if err != nil {
		return apps.Exception("failed to settle remote provider", err, zap.Any("", c), c.Logger)
	}
	c.Viper.SetConfigType("json")
	err = c.Viper.ReadRemoteConfig()
	if err != nil {
		return apps.Exception("failed to read remote config", err, zap.Any("", c), c.Logger)
	}
	return nil
}
