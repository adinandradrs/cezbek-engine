package handler

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/management"
	"github.com/gofiber/fiber/v2"
)

type PartnerManagement struct {
	management.PartnerManager
}

func newPartnerManagementResource(p PartnerManagement) *PartnerManagement {
	return &p
}

func PartnerManagementHandler(router fiber.Router, pm PartnerManagement) {
	handler := newPartnerManagementResource(pm)
	router.Post("/", handler.add)
}

// @Tags Partner Management APIs
// API Add Partner
// @Summary API Add Partner
// @Description API to register a new B2B Partner data as user and client
// @Schemes
// @Accept json
// @Param Authorization header string false "Your Token to Access" default(Bearer )
// @Param x-client-channel header string true "Client Channel" Enums(EBIZKEZBEK, B2BCLIENT)
// @Param x-client-os  header string true "Client OS or Browser Agent" default(android 10)
// @Param x-client-device  header string true "Client Device ID"
// @Param x-client-version  header string true "Client Platform Version" default(1.0.0)
// @Param x-client-timestamp  header string false "Client Original Timestamp in UNIX format (EPOCH)"
// @Param partner formData string true "Partner Corporate" default(PT. Lajada Piranti Commerce)
// @Param code formData string true "Partner Code" default(LAJADA)
// @Param email formData string true "Partner Email" default(kezbek.support@lajada.net)
// @Param msisdn formData string true "MSISDN" default(628123456789)
// @Param officer formData string true "Partner Officer" default(John Doe)
// @Param address formData string true "Office Address" default(Bintaro Exchange Mall Blok A1)
// @Param logo formData file true "Logo"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Failure 503
// @Router /v1/partners [post]
func (p *PartnerManagement) add(ctx *fiber.Ctx) error {
	logo, _ := ctx.FormFile("logo")
	inp := model.AddPartnerRequest{
		Partner: ctx.FormValue("partner"),
		Code:    ctx.FormValue("code"),
		Email:   ctx.FormValue("email"),
		Msisdn:  ctx.FormValue("msisdn"),
		Officer: ctx.FormValue("officer"),
		Address: ctx.FormValue("address"),
		Logo:    *logo,
		SessionRequest: model.SessionRequest{
			Id: 0,
		},
	}
	bad := apps.ValidateStruct(checker.Struct(inp))
	if bad != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(bad)
	}
	v, ex := p.Add(&inp)

	if ex != nil && ex.ErrorCode == apps.ErrCodeBussPartnerExists {
		return ctx.Status(fiber.StatusBadRequest).JSON(apps.BusinessErrorResponse(ex))
	}
	if ex != nil && ex.ErrorCode == apps.ErrCodeESBUnavailable {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(apps.BusinessErrorResponse(ex))
	}
	if ex != nil && (ex.ErrorCode == apps.ErrCodeSubmitted || ex.ErrorCode == apps.ErrCodeSomethingWrong) {
		return ctx.Status(fiber.StatusInternalServerError).JSON(apps.BusinessErrorResponse(ex))
	}

	return ctx.JSON(apps.DefaultSuccessResponse(apps.SuccessMsgSubmit, v))
}
