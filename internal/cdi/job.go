package cdi

import "github.com/adinandradrs/cezbek-engine/internal/usecase/job"

type JobUsecase struct {
	JobOnboardWatcher job.OnboardWatcher
}

func (c *Container) RegisterJobUsecase(infra Infra) JobUsecase {
	qNotificationEmailOtp := c.Viper.GetString("aws.sqs.topic.notification_email_otp")
	return JobUsecase{
		JobOnboardWatcher: job.NewOnboard(job.Onboard{
			Logger:                    c.Logger,
			QueueNotificationEmailOtp: &qNotificationEmailOtp,
			SqsAdapter:                infra.SQSAdapter,
			SesAdapter:                infra.SESAdapter,
		}),
	}
}
