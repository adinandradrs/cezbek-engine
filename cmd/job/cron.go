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
	job := gocron.NewScheduler(time.UTC)
	c.Logger.Info("cezbek cron job is running on background...")

	pool := goredis.NewPool(rv8.NewClient(
		&rv8.Options{
			Addr: c.Viper.GetString("rds.addresses"),
			DB:   c.Viper.GetInt("rds.index"),
		}))
	r := runner{
		Scheduler: job,
		Container: c,
		JobUsecase: c.RegisterJobUsecase(c.LoadInfra(),
			c.LoadRedis()),
		Redsync: redsync.New(pool),
	}
	r.onStartupJobExpireTier()
	r.onStartupJobSendInvoiceEmail()
	r.onStartupJobSendOtpEmail()
	job.StartBlocking()
}

type runner struct {
	*gocron.Scheduler
	cdi.Container
	cdi.JobUsecase
	*redsync.Redsync
}

func (r *runner) onStartupJobExpireTier() {
	_, err := r.Cron(r.Viper.GetString("schedule.expire_tier")).Do(func() {
		r.Logger.Info("expire_tier running...")
		mtx := r.NewMutex("expire_tier")
		if err := mtx.Lock(); err != nil {
			r.Logger.Error("expire_tier lock", zap.Error(err))
		}
		r.JobTierWatcher.Expire()
		if ok, err := mtx.Unlock(); !ok || err != nil {
			r.Logger.Error("expire_tier unlock", zap.Error(err))
		}
	})
	if err != nil {
		r.Logger.Panic("cezbek cron job is failing to run [JobTierWatcher.Expire]")
	}
}

func (r *runner) onStartupJobSendInvoiceEmail() {
	_, err := r.Every(r.Viper.GetString("schedule.send_invoice_email")).Do(func() {
		r.Logger.Info("send_invoice_email running...")
		mtx := r.NewMutex("send_invoice_email")
		if err := mtx.Lock(); err != nil {
			r.Logger.Error("send_invoice_email lock", zap.Error(err))
		}
		_ = r.JobTransactionWatcher.SendInvoiceEmail()
		if ok, err := mtx.Unlock(); !ok || err != nil {
			r.Logger.Error("send_invoice_email unlock", zap.Error(err))
		}
	})
	if err != nil {
		r.Logger.Panic("cezbek cron job is failing to run [JobTransactionWatcher.SendInvoiceEmail]")
	}
}

func (r *runner) onStartupJobSendOtpEmail() {
	_, err := r.Every(r.Viper.GetString("schedule.send_otp_email")).Do(func() {
		r.Logger.Info("send_otp_email running...")
		mtx := r.NewMutex("send_otp_email")
		if err := mtx.Lock(); err != nil {
			r.Logger.Error("send_otp_email lock", zap.Error(err))
		}
		_ = r.JobOnboardWatcher.SendOtpEmail()
		if ok, err := mtx.Unlock(); !ok || err != nil {
			r.Logger.Error("send_otp_email unlock", zap.Error(err))
		}
	})
	if err != nil {
		r.Logger.Panic("cezbek cron job is failing to run [JobOnboardWatcher.SendOtpEmail]")
	}
}
