package handler

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/handler/middleware"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/partner"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type PartnerTransaction struct {
	partner.TransactionProvider
	PartnerFilter fiber.Handler
}

func newPartnerTransaction(pt PartnerTransaction) *PartnerTransaction {
	return &pt
}

func PartnerTransactionHandler(router fiber.Router, pt PartnerTransaction) {
	handler := newPartnerTransaction(pt)
	router.Use(pt.PartnerFilter)
	router.Get("/", handler.search)
	router.Get("/:id", handler.detail)
}

// @Tags Transaction Partner APIs
// API Transaction Search
// @Summary API Transaction Search
// @Description API to search transaction by partner
// @Schemes
// @Accept json
// @Param Authorization header string true "Your Token to Access" default(Bearer )
// @Param x-client-channel header string true "Client Channel" Enums(EBIZKEZBEK, B2BCLIENT)
// @Param x-client-os  header string true "Client OS or Browser Agent" default(android 10)
// @Param x-client-device  header string true "Client Device ID"
// @Param x-client-version  header string true "Client Platform Version" default(1.0.0)
// @Param x-client-timestamp  header string false "Client Original Timestamp in UNIX format (EPOCH)"
// @Param Payload query model.SearchRequest true "Search Payload"
// @Success 200 {object} model.PartnerTransactionSearchResponse
// @Failure 400 {object} model.Meta
// @Failure 401 {object} model.Meta
// @Failure 403 {object} model.Meta
// @Failure 500 {object} model.Meta
// @Failure 503 {object} model.Meta
// @Router /partner/v1/transactions [get]
func (pt *PartnerTransaction) search(ctx *fiber.Ctx) error {
	inp := model.SearchRequest{}
	if err := ctx.QueryParser(&inp); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}
	inp.SessionRequest = middleware.ClientSession(ctx)
	v, ex := pt.Search(&inp)
	if ex != nil && ex.ErrorCode == apps.ErrCodeNotFound {
		return ctx.Status(fiber.StatusOK).
			JSON(apps.BusinessErrorResponse(ex))
	}
	return ctx.JSON(apps.DefaultSuccessResponse(apps.SuccessMsgDataFound, v))
}

// @Tags Transaction Partner APIs
// API Transaction Detail
// @Summary API Transaction Detail
// @Description API to view detail transaction by partner
// @Schemes
// @Accept json
// @Param Authorization header string true "Your Token to Access" default(Bearer )
// @Param x-client-channel header string true "Client Channel" Enums(EBIZKEZBEK, B2BCLIENT)
// @Param x-client-os  header string true "Client OS or Browser Agent" default(android 10)
// @Param x-client-device  header string true "Client Device ID"
// @Param x-client-version  header string true "Client Platform Version" default(1.0.0)
// @Param x-client-timestamp  header string false "Client Original Timestamp in UNIX format (EPOCH)"
// @Param id path int true "Transaction ID"
// @Success 200 {object} model.PartnerTransactionProjection
// @Failure 400 {object} model.Meta
// @Failure 401 {object} model.Meta
// @Failure 403 {object} model.Meta
// @Failure 500 {object} model.Meta
// @Failure 503 {object} model.Meta
// @Router /partner/v1/transactions/{id} [get]
func (pt *PartnerTransaction) detail(ctx *fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
	v, ex := pt.Detail(&model.FindByIdRequest{
		Id:             id,
		SessionRequest: middleware.ClientSession(ctx),
	})
	if ex != nil && ex.ErrorCode == apps.ErrCodeNotFound {
		return ctx.Status(fiber.StatusOK).
			JSON(apps.BusinessErrorResponse(ex))
	}
	return ctx.JSON(apps.DefaultSuccessResponse(apps.SuccessMsgDataFound, v))
}
