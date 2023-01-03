package handler

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/client"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/partner"
	"github.com/gofiber/fiber/v2"
)

type Authorization struct {
	PartnerOnboardProvider partner.OnboardProvider
	ClientOnboardProvider  client.OnboardProvider
	ClientFilter           fiber.Handler
}

func newAuthorizationResource(a Authorization) *Authorization {
	return &a
}

func AuthorizationHandler(router fiber.Router, auth Authorization) {
	handler := newAuthorizationResource(auth)
	router.Use(auth.ClientFilter).Post("/client", handler.clientAuth)
	router.Post("/b2b", handler.b2bAuth)
	router.Post("/otp", handler.otpAuth)
}

// @Tags Authorization APIs
// API Client Authorization
// @Summary API Client Authorization
// @Description API to authorize client's signature and code
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
	bad := apps.ValidateStruct(checker.Struct(inp))
	if bad != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(bad)
	}
	inp.ApiKey = ctx.Get(apps.HeaderApiKey)
	v, ex := a.ClientOnboardProvider.Authenticate(&inp)
	if ex != nil && ex.ErrorCode == apps.ErrCodeUnauthorized {
		return ctx.Status(fiber.StatusUnauthorized).JSON(apps.BusinessErrorResponse(ex))
	}
	if ex != nil && ex.ErrorCode == apps.ErrCodeSomethingWrong {
		return ctx.Status(fiber.StatusInternalServerError).JSON(apps.BusinessErrorResponse(ex))
	}
	return ctx.JSON(apps.DefaultSuccessResponse(apps.SuccessMsgSubmit, v))
}

// @Tags Authorization APIs
// API B2B Authorization
// @Summary API B2B Authorization
// @Description API to authorize B2B officer account
// @Schemes
// @Accept json
// @Param x-client-channel header string true "Client Channel" Enums(EBIZKEZBEK, B2BCLIENT)
// @Param x-client-os  header string true "Client OS or Browser Agent" default(android 10)
// @Param x-client-device  header string true "Client Device ID"
// @Param x-client-version  header string true "Client Platform Version" default(1.0.0)
// @Param x-client-timestamp  header string false "Client Original Timestamp in UNIX format (EPOCH)"
// @Param request body model.OfficerAuthenticationRequest true "B2B Officer Authentication Payload"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /v1/authorization/b2b [post]
func (a *Authorization) b2bAuth(ctx *fiber.Ctx) error {
	inp := model.OfficerAuthenticationRequest{}
	if err := ctx.BodyParser(&inp); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}
	bad := apps.ValidateStruct(checker.Struct(inp))
	if bad != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(bad)
	}
	v, ex := a.PartnerOnboardProvider.Authenticate(&inp)
	if ex != nil && ex.ErrorCode == apps.ErrCodeUnauthorized {
		return ctx.Status(fiber.StatusUnauthorized).JSON(apps.BusinessErrorResponse(ex))
	}
	if ex != nil && ex.ErrorCode == apps.ErrCodeSomethingWrong {
		return ctx.Status(fiber.StatusInternalServerError).JSON(apps.BusinessErrorResponse(ex))
	}
	return ctx.JSON(apps.DefaultSuccessResponse(apps.SuccessMsgSubmit, v))
}

// @Tags Authorization APIs
// API B2B Validation
// @Summary API B2B OTP Validation
// @Description API to validate B2B officer account OTP
// @Schemes
// @Accept json
// @Param x-client-channel header string true "Client Channel" Enums(EBIZKEZBEK, B2BCLIENT)
// @Param x-client-os  header string true "Client OS or Browser Agent" default(android 10)
// @Param x-client-device  header string true "Client Device ID"
// @Param x-client-version  header string true "Client Platform Version" default(1.0.0)
// @Param x-client-timestamp  header string false "Client Original Timestamp in UNIX format (EPOCH)"
// @Param x-client-trxid  header string true "Client Transaction ID"
// @Param request body model.OfficerValidationRequest true "B2B Officer Authentication Payload"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /v1/authorization/otp [post]
func (a *Authorization) otpAuth(ctx *fiber.Ctx) error {
	inp := model.OfficerValidationRequest{}
	if err := ctx.BodyParser(&inp); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}
	inp.TransactionId = ctx.Get(apps.HeaderClientTrxId)
	bad := apps.ValidateStruct(checker.Struct(inp))
	if bad != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(bad)
	}
	v, ex := a.PartnerOnboardProvider.Validate(&inp)
	if ex != nil && ex.ErrorCode == apps.ErrCodeBussPartnerOTPInvalid {
		return ctx.Status(fiber.StatusBadRequest).JSON(apps.BusinessErrorResponse(ex))
	}
	if ex != nil && ex.ErrorCode == apps.ErrCodeSomethingWrong {
		return ctx.Status(fiber.StatusInternalServerError).JSON(apps.BusinessErrorResponse(ex))
	}
	return ctx.JSON(apps.DefaultSuccessResponse(apps.SuccessMsgSubmit, v))
}
