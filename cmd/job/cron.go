package main

import (
	"github.com/adinandradrs/cezbek-engine/internal/cdi"
	"github.com/go-co-op/gocron"
	"time"
)

func main() {
	c := cdi.NewContainer("app_cezbek_job")
	infra := c.LoadInfra()
	redis := c.LoadRedis()
	ucase := c.RegisterJobUsecase(infra, redis)
	job := gocron.NewScheduler(time.UTC)
	c.Logger.Info("cezbek cron job is running on background...")
	_, err := job.Every(c.Viper.GetString("schedule.send_otp_email")).Do(func() {
		_ = ucase.JobOnboardWatcher.SendOtpEmail()
	})
	if err != nil {
		c.Logger.Panic("cezbek cron job is failing to run [JobOnboardManager.SendOtpEmail]")
	}

	job.StartBlocking()
}
