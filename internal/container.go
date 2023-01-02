package internal

import (
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/job"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/management"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/partner"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

type (
	Container struct {
		Logger     *zap.Logger
		HttpLogger fiber.Handler
		Viper      *viper.Viper
		app        string
	}

	Env struct {
		ContextPath string
		HttpPort    string
	}

	Infra struct {
		adaptor.S3Watcher
		adaptor.SQSAdapter
		adaptor.SESAdapter
		CiamPartner adaptor.CiamWatcher
	}

	Dao struct {
		repository.PartnerPersister
		repository.ParamPersister
		repository.H2HPersister
	}

	Usecase struct {
		management.PartnerManager
		management.ParamManager
		management.H2HManager
		PartnerOnboardManager partner.OnboardManager
		JobOnboardManager     job.OnboardManager
	}
)

func svcRegister(c *Container) Container {
	clientSvc := adaptor.NewConsul(adaptor.Consul{
		Host:    c.Viper.GetString("consul_host"),
		Port:    c.Viper.GetInt("consul_port"),
		Service: c.Viper.GetString(c.app),
		Viper:   c.Viper,
		Logger:  c.Logger,
	})
	ex := clientSvc.Register()
	if ex != nil {
		c.Logger.Panic("error to register", zap.Any("", &ex))
	}
	return *c
}

func NewContainer(app string) Container {
	decimal.MarshalJSONWithoutQuotes = true
	logger, httpLogger := apps.NewLog(false)
	conf, err := apps.NewEnv(logger)
	if err != nil {
		logger.Panic("error to load config", zap.Any("", &err))
	}
	return svcRegister(&Container{
		Logger:     logger,
		HttpLogger: httpLogger,
		Viper:      conf,
		app:        app,
	})
}

func (c *Container) loadPool() *storage.PgPool {
	opts := c.Viper.GetString("db_options")
	return storage.NewPgPool(&storage.PgOptions{
		Host:    c.Viper.GetString("db.host"),
		Port:    c.Viper.GetString("db.port"),
		User:    c.Viper.GetString("db.username"),
		Passwd:  c.Viper.GetString("db.password"),
		Schema:  c.Viper.GetString("db.schema"),
		Options: &opts,
		Logger:  c.Logger,
	})
}

func (c *Container) registerRepository() Dao {
	p := c.loadPool()
	return Dao{
		PartnerPersister: repository.NewPartner(repository.Partner{Logger: c.Logger, Pool: p.Pool}),
		ParamPersister:   repository.NewParameter(repository.Parameter{Logger: c.Logger, Pool: p.Pool}),
		H2HPersister:     repository.NewH2H(repository.H2H{Logger: c.Logger, Pool: p.Pool}),
	}
}

func (c *Container) LoadInfra() Infra {
	kid, skey := c.Viper.GetString("aws.keyid"), c.Viper.GetString("aws.keysecret")
	jwkb, _ := json.Marshal(c.Viper.Get("aws.ciam.partner.jwk"))
	jwk := string(jwkb)
	sender := c.Viper.GetString("aws.ses.sender")
	sjkt, _ := session.NewSession(&aws.Config{
		Region:      aws.String(c.Viper.GetString("aws.region.jkt")),
		Credentials: credentials.NewStaticCredentials(kid, skey, ""),
	})
	return Infra{
		S3Watcher: adaptor.NewS3(adaptor.S3Bucket{
			Bucket:   c.Viper.GetString("aws.s3.bucket"),
			Uploader: s3manager.NewUploader(sjkt),
			Logger:   c.Logger,
		}),
		SQSAdapter: adaptor.NewSQS(adaptor.SQS{
			SQS: sqs.New(sjkt),
		}),
		SESAdapter: adaptor.NewSES(adaptor.SES{
			SES:    sesv2.New(sjkt),
			Sender: &sender,
			Logger: c.Logger,
		}),
		CiamPartner: adaptor.NewCognito(adaptor.Cognito{
			Provider: c.LoadCognito(c.Viper.GetString("aws.region.sgp")),
			ClientId: c.Viper.GetString("aws.ciam.partner.clientid"),
			UserPool: c.Viper.GetString("aws.ciam.partner.poolid"),
			Scrt:     c.Viper.GetString("aws.ciam.partner.secret"),
			Region:   c.Viper.GetString("aws.ciam.region"),
			JWK:      jwk,
			Logger:   c.Logger,
		}),
	}
}

func (c *Container) LoadCognito(region string) *cognito.CognitoIdentityProvider {
	cfg := &aws.Config{Region: aws.String(region)}
	ssgp, err := session.NewSession(cfg)
	if err != nil {
		c.Logger.Panic("failed to load CIAM cognito", zap.Error(err))
	}
	return cognito.New(ssgp)
}

func (c *Container) LoadRedis() (rds storage.Cacher) {
	addrs := strings.Split(c.Viper.GetString("rds.addresses"), ";")
	pwd := c.Viper.GetString("rds.password")
	if c.Viper.GetBool("rds.cluster") {
		rds = storage.NewClusterRedis(&storage.RedisOptions{
			Addrs:  addrs,
			Passwd: pwd,
			Logger: c.Logger,
		})
	} else {
		rds = storage.NewRedis(&storage.RedisOptions{
			Addr:   addrs[0],
			Passwd: pwd,
			Index:  c.Viper.GetInt("rds_index"),
			Logger: c.Logger,
		})
	}
	return rds
}

func (c *Container) RegisterUsecase(infra Infra, cacher storage.Cacher) Usecase {
	dao := c.registerRepository()
	cdn := c.Viper.GetString("aws.cdn_base")
	path := c.Viper.GetString("aws.s3.path")
	otpTtl := c.Viper.GetDuration("ttl.otp")
	qNotificationEmailOtp := c.Viper.GetString("aws.sqs.topic.notification_email_otp")
	return Usecase{
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
		PartnerOnboardManager: partner.NewOnboard(partner.Onboard{
			Dao:                       dao.PartnerPersister,
			Cacher:                    cacher,
			SqsAdapter:                infra.SQSAdapter,
			ClientAuthTTL:             c.Viper.GetDuration("ttl.client_auth"),
			CDN:                       &cdn,
			OtpTTL:                    otpTtl,
			CiamWatcher:               infra.CiamPartner,
			QueueNotificationEmailOtp: &qNotificationEmailOtp,
			Logger:                    c.Logger,
		}),
		JobOnboardManager: job.NewOnboard(job.Onboard{
			Logger:                    c.Logger,
			QueueNotificationEmailOtp: &qNotificationEmailOtp,
			SqsAdapter:                infra.SQSAdapter,
			SesAdapter:                infra.SESAdapter,
		}),
		H2HManager: management.NewH2H(management.H2H{
			Logger: c.Logger,
			Dao:    dao.H2HPersister,
			Cacher: cacher,
		}),
	}
}

func (c *Container) LoadEnv() Env {
	return Env{
		ContextPath: c.Viper.GetString("base_path"),
		HttpPort:    c.Viper.GetString("app_port"),
	}
}
