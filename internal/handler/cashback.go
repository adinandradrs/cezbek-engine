package handler

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/handler/middleware"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/client"
	"github.com/gofiber/fiber/v2"
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
	router.Get("/:trxId", handler.detail)
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
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Failure 503
// @Router /v1/cashbacks [post]
func (c *Cashback) add(ctx *fiber.Ctx) error {
	inp := model.TransactionRequest{}
	if err := ctx.BodyParser(&inp); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}
	bad := apps.ValidateStruct(checker.Struct(inp))
	inp.SessionRequest = middleware.ClientSession(ctx)
	if bad != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(bad)
	}
	v, ex := c.Add(&inp)
	if ex != nil && ex.ErrorCode == apps.ErrCodeBussMerchantCodeInvalid {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(apps.BusinessErrorResponse(ex))
	}
	if ex != nil && ex.ErrorCode == apps.ErrCodeBussClientAddTransaction {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(apps.BusinessErrorResponse(ex))
	}
	return ctx.JSON(apps.DefaultSuccessResponse(apps.SuccessMsgSubmit, v))
}

// @Tags Client Cashback APIs
// API Detail Applied Cashback
// @Summary API Detail Applied Cashback
// @Description API to view detail of applied cashback based on the given Kezbek transaction reference
// @Schemes
// @Accept json
// @Param Authorization header string false "Your Token to Access" default(Bearer )
// @Param x-client-channel header string true "Client Channel" Enums(EBIZKEZBEK, B2BCLIENT)
// @Param x-client-os  header string true "Client OS or Browser Agent" default(android 10)
// @Param x-client-device  header string true "Client Device ID"
// @Param x-client-version  header string true "Client Platform Version" default(1.0.0)
// @Param x-client-timestamp  header string false "Client Original Timestamp in UNIX format (EPOCH)"
// @Param trxId path string true "Kezbek Transaction Reference"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Failure 503
// @Router /v1/cashbacks/{trxId} [post]
func (c *Cashback) detail(ctx *fiber.Ctx) error {
	return nil
}
