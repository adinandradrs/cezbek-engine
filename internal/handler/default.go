package handler

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var checker = validator.New()

func DefaultHandler(r fiber.Router, p string) {
	r.Get(p+"/ping", ping)
}

// Ping godoc
// @Summary Show the status of server.
// @Description Ping the status of server, should be respond fastly.
// @Tags Default
// @Accept */*
// @Produce json
// @Success 200
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /ping [get]
func ping(ctx *fiber.Ctx) error {
	return ctx.JSON(apps.DefaultSuccessResponse("succeeded ping with pong!", "pong"))
}
