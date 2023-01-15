package main

import (
	"github.com/adinandradrs/cezbek-engine/internal/cdi"
	"github.com/go-co-op/gocron"
	rv8 "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.uber.org/zap"
	"time"
)

func main() {
	c := cdi.NewContainer("app_cezbek_job")
	infra := c.LoadInfra()
	redis := c.LoadRedis()
	ucase := c.RegisterJobUsecase(infra, redis)
	job := gocron.NewScheduler(time.UTC)
	c.Logger.Info("cezbek cron job is running on background...")

	client := rv8.NewClient(&rv8.Options{
		Addr: c.Viper.GetString("rds.addresses"),
		DB:   c.Viper.GetInt("rds.index"),
	})
	rs := redsync.New(goredis.NewPool(client))
	_, err := job.Every(c.Viper.GetString("schedule.send_otp_email")).Do(func() {
		c.Logger.Info("send_otp_email running...")
		mtx := rs.NewMutex("send_otp_email")
		if err := mtx.Lock(); err != nil {
			c.Logger.Error("send_otp_email lock", zap.Error(err))
		}
		_ = ucase.JobOnboardWatcher.SendOtpEmail()
		if ok, err := mtx.Unlock(); !ok || err != nil {
			c.Logger.Error("send_otp_email unlock", zap.Error(err))
		}
	})
	if err != nil {
		c.Logger.Panic("cezbek cron job is failing to run [JobOnboardWatcher.SendOtpEmail]")
	}

	_, err = job.Every(c.Viper.GetString("schedule.send_invoice_email")).Do(func() {
		c.Logger.Info("send_invoice_email running...")
		mtx := rs.NewMutex("send_invoice_email")
		if err := mtx.Lock(); err != nil {
			c.Logger.Error("send_invoice_email lock", zap.Error(err))
		}
		_ = ucase.JobTransactionWatcher.SendInvoiceEmail()
		if ok, err := mtx.Unlock(); !ok || err != nil {
			c.Logger.Error("send_invoice_email unlock", zap.Error(err))
		}
	})
	if err != nil {
		c.Logger.Panic("cezbek cron job is failing to run [JobTransactionWatcher.SendInvoiceEmail]")
	}

	_, err = job.Cron(c.Viper.GetString("schedule.expire_tier")).Do(func() {
		c.Logger.Info("expire_tier running...")
		mtx := rs.NewMutex("expire_tier")
		if err := mtx.Lock(); err != nil {
			c.Logger.Error("expire_tier lock", zap.Error(err))
		}
		ucase.JobTierWatcher.Expire()
		if ok, err := mtx.Unlock(); !ok || err != nil {
			c.Logger.Error("expire_tier unlock", zap.Error(err))
		}
	})
	if err != nil {
		c.Logger.Panic("cezbek cron job is failing to run [JobTierWatcher.Expire]")
	}

	job.StartBlocking()
}
