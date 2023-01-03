package cdi

import (
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/client"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/management"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/partner"
)

func (c *Container) RegisterAPIUsecase(infra Infra, cacher storage.Cacher) APIUsecase {
	dao := c.registerRepository()
	cdn := c.Viper.GetString("aws.cdn_base")
	path := c.Viper.GetString("aws.s3.path")
	otpTtl := c.Viper.GetDuration("ttl.otp")
	qNotificationEmailOtp := c.Viper.GetString("aws.sqs.topic.notification_email_otp")
	return APIUsecase{
		PartnerManager: management.NewPartner(management.Partner{
			Dao:         dao.PartnerPersister,
			CiamWatcher: infra.CiamPartner,
			S3Watcher:   infra.S3Watcher,
			PathS3:      &path,
			Logger:      c.Logger,
		}),
		ParamManager: management.NewParameter(management.Parameter{
			Dao:    dao.ParamPersister,
			Cacher: cacher,
			Logger: c.Logger,
		}),
		PartnerOnboardProvider: partner.NewOnboard(partner.Onboard{
			Dao:                       dao.PartnerPersister,
			Cacher:                    cacher,
			SqsAdapter:                infra.SQSAdapter,
			AuthTTL:                   c.Viper.GetDuration("ttl.client_auth"),
			CDN:                       &cdn,
			OtpTTL:                    otpTtl,
			CiamWatcher:               infra.CiamPartner,
			QueueNotificationEmailOtp: &qNotificationEmailOtp,
			Logger:                    c.Logger,
		}),
		ClientOnboardProvider: client.NewOnboard(client.Onboard{
			Dao:         dao.PartnerPersister,
			Cacher:      cacher,
			AuthTTL:     c.Viper.GetDuration("ttl.client_auth"),
			CiamWatcher: infra.CiamPartner,
			Logger:      c.Logger,
		}),
		H2HManager: management.NewH2H(management.H2H{
			Logger: c.Logger,
			Dao:    dao.H2HPersister,
			Cacher: cacher,
		}),
	}
}
