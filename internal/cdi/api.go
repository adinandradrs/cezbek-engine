package cdi

import (
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/client"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/h2h"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/management"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/partner"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/workflow"
)

type APIUsecase struct {
	management.PartnerManager
	management.ParamManager
	management.H2HManager
	management.WorkflowManager
	workflow.CashbackProvider
	PartnerOnboardProvider    partner.OnboardProvider
	ClientOnboardProvider     client.OnboardProvider
	ClientTransactionProvider client.TransactionProvider
	H2HFactory                h2h.Factory
}

func (c *Container) RegisterAPIUsecase(infra Infra, cacher storage.Cacher) APIUsecase {
	dao := c.registerRepository()
	cdn := c.Viper.GetString("aws.cdn_base")
	path := c.Viper.GetString("aws.s3.path")
	otpTtl := c.Viper.GetDuration("ttl.otp")
	qNotificationEmailOtp := c.Viper.GetString("aws.sqs.topic.notification_email_otp")
	qNotificationEmailInvoice := c.Viper.GetString("aws.sqs.topic.notification_email_invoice")
	tierProvider := workflow.NewTier(workflow.Tier{
		Dao:            dao.TierPersister,
		Logger:         c.Logger,
		Cacher:         cacher,
		ExpiryDuration: c.Viper.GetDuration("wfreward.expiry_duration"),
	})
	h2hFactory := h2h.NewFactory(h2h.Factory{
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
	})
	cashbackProvider := workflow.NewCashback(workflow.Cashback{
		Logger: c.Logger,
		Dao:    dao.WorkflowPersister,
	})
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
		CashbackProvider: workflow.NewCashback(workflow.Cashback{
			Logger: c.Logger,
			Dao:    dao.WorkflowPersister,
		}),
		H2HManager: management.NewH2H(management.H2H{
			Logger: c.Logger,
			Dao:    dao.H2HPersister,
			Cacher: cacher,
		}),
		WorkflowManager: management.NewWorkflow(management.Workflow{
			Logger: c.Logger,
			Dao:    dao.WorkflowPersister,
			Cacher: cacher,
		}),
		ClientTransactionProvider: client.NewTransaction(client.Transaction{
			Dao:                           dao.TransactionPersister,
			CashbackProvider:              cashbackProvider,
			QueueNotificationEmailInvoice: &qNotificationEmailInvoice,
			TierProvider:                  tierProvider,
			SqsAdapter:                    infra.SQSAdapter,
			Factory:                       h2hFactory,
			Logger:                        c.Logger,
			Cacher:                        cacher,
		}),
	}
}
