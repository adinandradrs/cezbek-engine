package adaptor

import (
	"fmt"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Consul struct {
	Host          string
	Port          int
	Service       string
	CheckInterval string
	CheckTimeout  string
	Logger        *zap.Logger
	Config        *viper.Viper
}

type ConsulWatcher interface {
	Register() *model.TechnicalError
}

func NewConsul(c Consul) ConsulWatcher {
	return &c
}

func (c Consul) Register() *model.TechnicalError {
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		c.Logger.Fatal("cannot initialize consul client", zap.Error(err))
	}
	reg := new(consul.AgentServiceRegistration)
	reg.ID = c.Service
	reg.Name = c.Service
	reg.Address = c.Host
	reg.Port = c.Port
	reg.Check = new(consul.AgentServiceCheck)
	reg.Check.Interval = c.CheckInterval
	reg.Check.Timeout = c.CheckTimeout
	err = client.Agent().ServiceRegister(reg)
	if err != nil {
		return &model.TechnicalError{
			Exception: err.Error(),
		}
	}
	c.readConfig()
	return nil
}

func (c Consul) readConfig() {
	err := c.Config.AddRemoteProvider("consul", fmt.Sprintf("%s:%v", c.Host, c.Port), c.Service)
	if err != nil {
		c.Logger.Fatal("viper cannot settle with remote consul", zap.Error(err))
	}
	c.Config.SetConfigType("json")
	err = c.Config.ReadRemoteConfig()
	if err != nil {
		c.Logger.Fatal("viper cannot read KV on remote consul", zap.Error(err))
	}
}
