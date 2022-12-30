package handler

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/management"
	"github.com/gofiber/fiber/v2"
)

type Partner struct {
	management.PartnerManager
}

func newPartnerResource(p Partner) *Partner {
	return &p
}

func PartnerHandler(r fiber.Router, p Partner) {
	handler := newPartnerResource(p)
	route := "/partners"
	r.Post(route, handler.add)
}

// @Tags Partner Management APIs
// Add Partner API
// @Summary Add Partner API
// @Schemes
// @Accept json
// @Param Authorization header string false "Your Token to Access" default(Bearer )
// @Param x-client-channel header string true "Client Channel" Enums(EBIZKEZBEK)
// @Param x-client-os  header string true "Client OS or Browser Agent" default(android 10)
// @Param x-client-device  header string true "Client Device ID"
// @Param x-client-version  header string true "Client Platform Version" default(1.0.0)
// @Param x-client-timestamp  header string false "Client Original Time Request in UNIX Timestamp"
// @Param partner formData string true "Partner Corporate" default(PT. Lajada Piranti Commerce)
// @Param code formData string true "Partner Code" default(LAJADA)
// @Param email formData string true "Partner Email" default(kezbek.support@lajada.net)
// @Param msisdn formData string true "MSISDN" default(628123456789)
// @Param officer formData string true "Partner Officer" default(John Doe)
// @Param address formData string true "Office Address" default(100000)
// @Param logo formData file true "Logo"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /v1/partners [post]
func (p *Partner) add(ctx *fiber.Ctx) error {
	v, ex := p.Add(nil)
	if ex != nil && ex.ErrorCode == apps.ErrCodeBussInvalidCodePartner {
		return ctx.Status(fiber.StatusBadRequest).JSON(ex)
	}
	return ctx.JSON(apps.DefaultSuccessResponse(apps.SuccessMsgSubmit, v))
}
