package main

import (
	"context"
	"github.com/adinandradrs/cezbek-engine/internal/cdi"
	"github.com/bsm/redislock"
	"github.com/go-co-op/gocron"
	rv9 "github.com/go-redis/redis/v9"
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

	ctx := context.Background()
	client := rv9.NewClient(&rv9.Options{
		Addr: c.Viper.GetString("rds.addresses"),
		DB:   c.Viper.GetInt("rds.index"),
	})
	locker := redislock.New(client)

	_, err := job.Every(c.Viper.GetString("schedule.send_otp_email")).Do(func() {
		lock, err := locker.Obtain(ctx, "send_otp_email", c.Viper.GetDuration("lock.send_otp_email"), nil)
		if err == redislock.ErrNotObtained {
			c.Logger.Warn("lock obtain failed send_otp_email", zap.String("", "lock.send_otp_email"), zap.Error(err))
		} else if err != nil {
			c.Logger.Panic("schedule.send_otp_email error lock", zap.Error(err))
		}
		_ = ucase.JobOnboardWatcher.SendOtpEmail()
		defer lock.Release(ctx)
	})
	if err != nil {
		c.Logger.Panic("cezbek cron job is failing to run [JobOnboardWatcher.SendOtpEmail]")
	}

	_, err = job.Every(c.Viper.GetString("schedule.send_invoice_email")).Do(func() {
		lock, err := locker.Obtain(ctx, "send_invoice_email", c.Viper.GetDuration("lock.send_invoice_email"), nil)
		if err == redislock.ErrNotObtained {
			c.Logger.Warn("lock obtain failed send_invoice_email", zap.String("", "lock.send_invoice_email"), zap.Error(err))
		} else if err != nil {
			c.Logger.Panic("schedule.send_invoice_email error lock", zap.Error(err))
		}
		_ = ucase.JobTransactionWatcher.SendInvoiceEmail()
		defer lock.Release(ctx)
	})
	if err != nil {
		c.Logger.Panic("cezbek cron job is failing to run [JobTransactionWatcher.SendInvoiceEmail]")
	}

	_, err = job.Cron(c.Viper.GetString("schedule.expire_tier")).Do(func() {
		lock, err := locker.Obtain(ctx, "expire_tier", c.Viper.GetDuration("lock.expire_tier"), nil)
		if err == redislock.ErrNotObtained {
			c.Logger.Warn("lock obtain failed expire_tier", zap.String("", "lock.expire_tier"), zap.Error(err))
		} else if err != nil {
			c.Logger.Panic("schedule.expire_tier error lock", zap.Error(err))
		}
		ucase.JobTierWatcher.Expire()
		defer lock.Release(ctx)
	})
	if err != nil {
		c.Logger.Panic("cezbek cron job is failing to run [JobTierWatcher.Expire]")
	}

	job.StartBlocking()
}
