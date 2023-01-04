package cdi

import (
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/h2h"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/job"
)

type JobUsecase struct {
	JobOnboardWatcher job.OnboardWatcher
	H2HFactory        h2h.Factory
}

func (c *Container) RegisterJobUsecase(infra Infra, cacher storage.Cacher) JobUsecase {
	qNotificationEmailOtp := c.Viper.GetString("aws.sqs.topic.notification_email_otp")
	return JobUsecase{
		JobOnboardWatcher: job.NewOnboard(job.Onboard{
			Logger:                    c.Logger,
			QueueNotificationEmailOtp: &qNotificationEmailOtp,
			SqsAdapter:                infra.SQSAdapter,
			SesAdapter:                infra.SESAdapter,
		}),
		H2HFactory: h2h.NewFactory(h2h.Factory{
			Cacher: cacher,
			Gopaid: h2h.Gopaid{GopaidAdapter: infra.GopaidAdapter},
			Josvo:  h2h.Josvo{JosvoAdapter: infra.JosvoAdapter},
			Linksaja: h2h.Linksaja{
				TokenTTL:        c.Viper.GetDuration("ttl.lsaja"),
				Cacher:          cacher,
				LinksajaAdapter: infra.LinksajaAdapter,
			},
			Xenit:       h2h.Xenit{XenitAdapter: infra.XenitAdapter},
			Middletrans: h2h.Middletrans{MiddletransAdapter: infra.MiddletransAdapter},
		}),
	}
}
