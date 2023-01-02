package main

import (
	"github.com/adinandradrs/cezbek-engine/internal"
	"github.com/go-co-op/gocron"
	"time"
)

func main() {
	c := internal.NewContainer("app_cezbek_job")
	infra := c.LoadInfra()
	redis := c.LoadRedis()
	ucase := c.RegisterUsecase(infra, redis)
	s := gocron.NewScheduler(time.UTC)
	c.Logger.Info("cezbek cron job is running on background...")
	_, err := s.Every(c.Viper.GetString("schedule.send_otp_email")).Do(func() {
		_ = ucase.JobOnboardManager.SendOtpEmail()
	})
	if err != nil {
		c.Logger.Panic("cezbek cron job is failing to run [JobOnboardManager.SendOtpEmail]")
	}

	s.StartBlocking()
}