package main

import (
	"github.com/adinandradrs/cezbek-engine/internal/cdi"
	_ "github.com/adinandradrs/cezbek-engine/internal/docs"
	"github.com/adinandradrs/cezbek-engine/internal/handler"
	"github.com/adinandradrs/cezbek-engine/internal/handler/middleware"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/management"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	swagger "github.com/swaggo/fiber-swagger"
	"go.uber.org/zap"
)

// @title Kezbek - Cashback Engine Sandbox
// @version 1.0-Beta
// @description This Cashback Engine Sandbox is only used for test and development purpose. To explore and serve all Kezbek operational APIs as a live data. It is not intended for production usage.
// @termsOfService http://swagger.io/terms/

// @contact.name Kezbek Developer
// @contact.url https://kezbek.id
// @contact.email developer@kezbek.id

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api
func main() {
	//load Contexts and Dependency Injection (CDI)
	c := cdi.NewContainer("app_cezbek_api")
	env := c.LoadEnv()
	infra := c.LoadInfra()
	redis := c.LoadRedis()
	ucase := c.RegisterAPIUsecase(infra, redis)

	//app event(s)
	onStartupLoadParameterCache(ucase.ParamManager, c.Logger)
	onStartupLoadH2HCache(ucase.H2HManager, c.Logger)
	onStartupLoadWorkflowCache(ucase.WorkflowManager, c.Logger)

	//app starting
	api := fiber.New()

	//middleware config
	api.Use(cors.New())
	preAuthenticator := middleware.NewPreAuthenticator(&middleware.PreAuthenticator{
		Logger: c.Logger,
	})
	preAuthClientFilter := preAuthenticator.ClientFilter()
	jwtAuthenticator := middleware.NewJwtAuthenticator(&middleware.JwtAuthenticator{
		Logger:      c.Logger,
		CiamPartner: infra.CiamPartner,
		Cacher:      redis,
	})
	jwtAuthClientFilter := jwtAuthenticator.ClientFilter()
	jwtAuthPartnerFilter := jwtAuthenticator.PartnerFilter()

	//swagger
	api.Get(env.ContextPath+"/swagger/*", swagger.WrapHandler)
	handler.DefaultHandler(api, env.ContextPath)

	//APIs
	authorization := api.Group("/api/v1/authorization").Use(c.HttpLogger)
	handler.AuthorizationHandler(authorization, handler.Authorization{
		PartnerOnboardProvider: ucase.PartnerOnboardProvider,
		ClientOnboardProvider:  ucase.ClientOnboardProvider,
		ClientFilter:           preAuthClientFilter,
	})

	partners := api.Group("/api/v1/partners").Use(c.HttpLogger)
	handler.PartnerManagementHandler(partners, handler.PartnerManagement{
		PartnerManager: ucase.PartnerManager,
	})

	cashbacks := api.Group("/api/v1/cashbacks").Use(c.HttpLogger)
	handler.CashbackHandler(cashbacks, handler.Cashback{
		TransactionProvider: ucase.ClientTransactionProvider,
		ClientFilter:        jwtAuthClientFilter,
	})

	partnerTransactions := api.Group("/api/partner/v1/transactions")
	handler.PartnerTransactionHandler(partnerTransactions, handler.PartnerTransaction{
		TransactionProvider: ucase.PartnerTransactionProvider,
		PartnerFilter:       jwtAuthPartnerFilter,
	})

	_ = api.Listen(env.HttpPort)
}

func onStartupLoadWorkflowCache(h management.WorkflowManager, logger *zap.Logger) {
	go func() {
		ex := h.CacheRewardTiers()
		if ex != nil {
			logger.Panic("failed to load reward tiers")
		}
	}()
}

func onStartupLoadH2HCache(h management.H2HManager, logger *zap.Logger) {
	go func() {
		ex := h.CacheProviders()
		if ex != nil {
			logger.Panic("failed to load h2h providers")
		}
	}()

	go func() {
		ex := h.CachePricelists()
		if ex != nil {
			logger.Panic("failed to load h2h pricelists")
		}
	}()
}

func onStartupLoadParameterCache(p management.ParamManager, logger *zap.Logger) {
	go func() {
		ex := p.CacheEmailSubjects()
		if ex != nil {
			logger.Panic("failed to load email subjects")
		}
	}()

	go func() {
		ex := p.CacheEmailTemplates()
		if ex != nil {
			logger.Panic("failed to load email templates")
		}
	}()

	go func() {
		ex := p.CacheWallets()
		if ex != nil {
			logger.Panic("failed to load wallet codes")
		}
	}()
}
