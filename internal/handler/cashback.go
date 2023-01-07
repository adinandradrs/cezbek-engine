package handler

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/handler/middleware"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/client"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

type Cashback struct {
	client.TransactionProvider
	ClientFilter fiber.Handler
}

func newCashback(c Cashback) *Cashback {
	return &c
}

func CashbackHandler(router fiber.Router, c Cashback) {
	handler := newCashback(c)
	router.Use(c.ClientFilter)
	router.Post("/", handler.add)
	router.Get("/:msisdn", handler.info)
}

// @Tags Client Cashback APIs
// API Apply Cashback
// @Summary API Apply Cashback
// @Description API to apply cashback on client's transaction
// @Schemes
// @Accept json
// @Param Authorization header string true "Your Token to Access" default(Bearer )
// @Param x-client-channel header string true "Client Channel" Enums(EBIZKEZBEK, B2BCLIENT)
// @Param x-client-os  header string true "Client OS or Browser Agent" default(android 10)
// @Param x-client-device  header string true "Client Device ID"
// @Param x-client-version  header string true "Client Platform Version" default(1.0.0)
// @Param x-client-timestamp  header string false "Client Original Timestamp in UNIX format (EPOCH)"
// @Param request body model.TransactionRequest true "Transaction Payload"
// @Success 200 {object} model.TransactionResponse
// @Failure 400 {object} model.Meta
// @Failure 401 {object} model.Meta
// @Failure 403 {object} model.Meta
// @Failure 500 {object} model.Meta
// @Failure 503 {object} model.Meta
// @Router /v1/cashbacks [post]
func (c *Cashback) add(ctx *fiber.Ctx) error {
	inp := model.TransactionRequest{}
	if err := ctx.BodyParser(&inp); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}
	bad := apps.ValidateStruct(checker.Struct(inp))
	if inp.Amount.Cmp(decimal.Zero) <= 0 {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(apps.BusinessErrorResponse(&model.BusinessError{
				ErrorCode:    apps.ErrCodeBadPayload,
				ErrorMessage: apps.ErrMsgBadPayload,
			}))
	}
	inp.SessionRequest = middleware.ClientSession(ctx)
	if bad != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(bad)
	}
	v, ex := c.Add(&inp)
	if ex != nil && (ex.ErrorCode == apps.ErrCodeBussMerchantCodeInvalid ||
		ex.ErrorCode == apps.ErrCodeBussNoCashback) {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(apps.BusinessErrorResponse(ex))
	}
	if ex != nil && (ex.ErrorCode == apps.ErrCodeBussH2HCashbackFailed ||
		ex.ErrorCode == apps.ErrCodeSomethingWrong ||
		ex.ErrorCode == apps.ErrCodeBussClientAddTransaction) {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(apps.BusinessErrorResponse(ex))
	}
	return ctx.JSON(apps.DefaultSuccessResponse(apps.SuccessMsgSubmit, v))
}

// @Tags Client Cashback APIs
// API Tier Information
// @Summary API Tier Information
// @Description API to retrieve tier information
// @Schemes
// @Accept json
// @Param Authorization header string true "Your Token to Access" default(Bearer )
// @Param x-client-channel header string true "Client Channel" Enums(EBIZKEZBEK, B2BCLIENT)
// @Param x-client-os  header string true "Client OS or Browser Agent" default(android 10)
// @Param x-client-device  header string true "Client Device ID"
// @Param x-client-version  header string true "Client Platform Version" default(1.0.0)
// @Param x-client-timestamp  header string false "Client Original Timestamp in UNIX format (EPOCH)"
// @Param msisdn path string true "Customer MSISDN"
// @Success 200 {object} model.TransactionTierResponse
// @Failure 400 {object} model.Meta
// @Failure 401 {object} model.Meta
// @Failure 403 {object} model.Meta
// @Failure 500 {object} model.Meta
// @Failure 503 {object} model.Meta
// @Router /v1/cashbacks/{msisdn} [get]
func (c *Cashback) info(ctx *fiber.Ctx) error {
	inp := middleware.ClientSession(ctx)
	inp.Msisdn = ctx.Params("msisdn")
	v, ex := c.Tier(&inp)
	if ex != nil && ex.ErrorCode == apps.ErrCodeNotFound {
		return ctx.Status(fiber.StatusOK).
			JSON(apps.BusinessErrorResponse(ex))
	}
	return ctx.JSON(apps.DefaultSuccessResponse(apps.SuccessMsgDataFound, v))
}
