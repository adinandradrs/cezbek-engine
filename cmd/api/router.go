package main

import (
	"github.com/adinandradrs/cezbek-engine/internal"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	_ "github.com/adinandradrs/cezbek-engine/internal/docs"
	"github.com/adinandradrs/cezbek-engine/internal/handler"
	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"
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

	app := fiber.New()
	app.Get(env.ContextPath+"/swagger/*", swagger.WrapHandler)
	handler.DefaultHandler(app, env.ContextPath)

	m := apps.Middleware{Logger: c.Logger}
	authenticator := apps.Authenticator(m)

	authorization := app.Group("/api/v1/authorization").Use(authenticator)
	handler.AuthorizationHandler(authorization, handler.Authorization{
		OnboardManager: ucase.OnboardManager,
	})

	partners := app.Group("/api/v1/partners")
	handler.PartnerHandler(partners, handler.Partner{
		PartnerManager: ucase.PartnerManager,
	})

	_ = app.Listen(env.HttpPort)
}
