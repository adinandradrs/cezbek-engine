package handler

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/partner"
	"github.com/gofiber/fiber/v2"
)

type Authorization struct {
	partner.OnboardManager
}

func newAuthorizationResource(a Authorization) *Authorization {
	return &a
}

func AuthorizationHandler(r fiber.Router, a Authorization) {
	h := newAuthorizationResource(a)
	r.Post("/client", h.clientAuth)
}

// @Tags Authorization APIs
// Client Authorization API
// @Summary Client Authorization API
// @Description This API is to authorize client's signature and code
// @Schemes
// @Accept json
// @Param x-client-signature header string true "Client signature using HMAC SHA256, signature formula is <b>HEX(HMAC(SHA256(UPPER(HTTP-METHOD):UPPER(CODE):UNIX-EPOCH:UPPER(API-KEY))))</b>"
// @Param x-api-key header string true "Client API Key"
// @Param x-client-channel header string true "Client Channel" Enums(EBIZKEZBEK, B2BCLIENT)
// @Param x-client-os  header string true "Client OS or Browser Agent" default(android 10)
// @Param x-client-device  header string true "Client Device ID"
// @Param x-client-version  header string true "Client Platform Version" default(1.0.0)
// @Param x-client-timestamp  header string true "Client Original Timestamp in UNIX format (EPOCH)"
// @Param request body model.ClientAuthenticationRequest true "Client Authentication Payload"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /v1/authorization/client [post]
func (a *Authorization) clientAuth(ctx *fiber.Ctx) error {
	inp := model.ClientAuthenticationRequest{}
	if err := ctx.BodyParser(&inp); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}
	inp.ApiKey = ctx.Get(apps.HeaderApiKey)
	v, ex := a.AuthenticateClient(&inp)
	if ex != nil && ex.ErrorCode == apps.ErrCodeUnauthorized {
		return ctx.Status(fiber.StatusUnauthorized).JSON(apps.BusinessErrorResponse(ex))
	}
	if ex != nil && ex.ErrorCode == apps.ErrCodeSomethingWrong {
		return ctx.Status(fiber.StatusInternalServerError).JSON(apps.BusinessErrorResponse(ex))
	}
	return ctx.JSON(apps.DefaultSuccessResponse(apps.SuccessMsgSubmit, v))
}
