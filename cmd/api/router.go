package main

import (
	"github.com/adinandradrs/cezbek-engine/internal"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	_ "github.com/adinandradrs/cezbek-engine/internal/docs"
	"github.com/adinandradrs/cezbek-engine/internal/handler"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/management"
	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"
	"go.uber.org/zap"
)

// @title Kezbek - Cashback Engine Sandbox
// @version 1.0
// @description This Cashback Engine Sandbox is only used for test and development purpose. To explore and serve all Kezbek operational APIs as a live data. It is not intended for production usage.
// @termsOfService http://swagger.io/terms/

// @contact.name Kezbek Developer
// @contact.url https://kezbek.id
// @contact.email developer@kezbek.id

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api
func main() {
	c := internal.NewContainer("app_cezbek_api")
	env := c.LoadEnv()
	infra := c.LoadInfra()
	redis := c.LoadRedis()
	ucase := c.RegisterUsecase(infra, redis)

	m := apps.Middleware{Logger: c.Logger}
	authenticator := apps.Authenticator(m)

	app := fiber.New()
	app.Get(env.ContextPath+"/swagger/*", swagger.WrapHandler)
	handler.DefaultHandler(app, env.ContextPath)

	authorization := app.Group("/api/v1/authorization").Use(c.HttpLogger, authenticator)
	handler.AuthorizationHandler(authorization, handler.Authorization{
		OnboardManager: ucase.PartnerOnboardManager,
	})

	partners := app.Group("/api/v1/partners").Use(c.HttpLogger)
	handler.PartnerHandler(partners, handler.Partner{
		PartnerManager: ucase.PartnerManager,
	})

	loadParameterCache(ucase.ParamManager, c.Logger)
	loadH2HCache(ucase.H2HManager, c.Logger)
	_ = app.Listen(env.HttpPort)
}

func loadH2HCache(h management.H2HManager, logger *zap.Logger) {
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

func loadParameterCache(p management.ParamManager, logger *zap.Logger) {
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
