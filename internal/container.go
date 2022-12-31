package internal

import (
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/management"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/partner"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

type (
	Container struct {
		Logger *zap.Logger
		Viper  *viper.Viper
		app    string
	}

	Env struct {
		ContextPath string
		HttpPort    string
	}

	Infra struct {
		adaptor.S3Watcher
		CiamPartner adaptor.CiamWatcher
	}

	Dao struct {
		repository.PartnerPersister
	}

	Usecase struct {
		management.PartnerManager
		partner.OnboardManager
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
	logger := apps.NewLog(false)
	conf, err := apps.NewEnv(logger)
	if err != nil {
		logger.Panic("error to load config", zap.Any("", &err))
	}
	return svcRegister(&Container{
		Logger: logger,
		Viper:  conf,
		app:    app,
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
	}
}

func (c *Container) LoadInfra() Infra {
	kid, skey := c.Viper.GetString("aws.keyid"), c.Viper.GetString("aws.keysecret")
	s3Credential, _ := session.NewSession(&aws.Config{
		Region:      aws.String(c.Viper.GetString("aws.s3.region")),
		Credentials: credentials.NewStaticCredentials(kid, skey, ""),
	})
	s3Watcher := adaptor.NewS3(adaptor.S3Bucket{
		Bucket:   c.Viper.GetString("aws.s3.bucket"),
		Uploader: s3manager.NewUploader(s3Credential),
		Logger:   c.Logger,
	})
	jwkb, _ := json.Marshal(c.Viper.Get("aws.ciam.partner.jwk"))
	jwk := string(jwkb)
	ciamPartner := adaptor.NewCognito(adaptor.Cognito{
		Provider: c.LoadCognito(c.Viper.GetString("aws.ciam.region")),
		ClientId: c.Viper.GetString("aws.ciam.partner.clientid"),
		UserPool: c.Viper.GetString("aws.ciam.partner.poolid"),
		Scrt:     c.Viper.GetString("aws.ciam.partner.secret"),
		Region:   c.Viper.GetString("aws.ciam.region"),
		JWK:      jwk,
		Logger:   c.Logger,
	})
	return Infra{
		S3Watcher:   s3Watcher,
		CiamPartner: ciamPartner,
	}
}

func (c *Container) LoadCognito(region string) *cognito.CognitoIdentityProvider {
	cfg := &aws.Config{Region: aws.String(region)}
	s, err := session.NewSession(cfg)
	if err != nil {
		c.Logger.Panic("failed to load CIAM cognito", zap.Error(err))
	}
	return cognito.New(s)
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

func (c *Container) RegisterUsecase(infra Infra, rds storage.Cacher) Usecase {
	dao := c.registerRepository()
	cdn := c.Viper.GetString("aws.cdn_base")
	path := c.Viper.GetString("aws.s3.path")
	return Usecase{
		PartnerManager: management.NewPartner(management.Partner{
			Dao:         dao.PartnerPersister,
			CiamWatcher: infra.CiamPartner,
			S3Watcher:   infra.S3Watcher,
			CDN:         cdn,
			PathS3:      path,
			Logger:      c.Logger,
		}),
		OnboardManager: partner.NewOnboard(partner.Onboard{
			Dao:           dao.PartnerPersister,
			Cacher:        rds,
			ClientAuthTTL: c.Viper.GetDuration("ttl.client_auth"),
			CiamWatcher:   infra.CiamPartner,
			Logger:        c.Logger,
		}),
	}
}

func (c *Container) LoadEnv() Env {
	return Env{
		ContextPath: c.Viper.GetString("base_path"),
		HttpPort:    c.Viper.GetString("app_port"),
	}
}
