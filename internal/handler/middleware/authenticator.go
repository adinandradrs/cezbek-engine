package middleware

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strings"
)

type Authenticator struct {
	Logger *zap.Logger
}

func NewAuthenticator(a *Authenticator) Authenticator {
	return *a
}

func (a *Authenticator) ClientFilter() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Get(apps.HeaderClientChannel) != apps.ChannelB2BClient {
			a.Logger.Error("the given channel is not valid", zap.String("signature", ctx.Get(apps.HeaderClientSignature)),
				zap.String("channel", ctx.Get(apps.HeaderClientChannel)))
			return ctx.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Meta: model.Meta{
					Code:    apps.ErrCodeInvalidChannel,
					Message: apps.ErrMsgInvalidChannel,
				},
			})
		}
		b := struct {
			Code string `json:"code"`
		}{}
		if err := ctx.BodyParser(&b); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err)
		}
		epoch := ctx.Get(apps.HeaderClientTimestamp)
		key := ctx.Get(apps.HeaderApiKey)

		d := string(ctx.Request().Header.Method()) + ":" + strings.ToUpper(b.Code) + ":" + epoch + ":" + strings.ToUpper(key)
		if apps.HMAC(d, b.Code) != ctx.Get(apps.HeaderClientSignature) {
			a.Logger.Error("the given signature is not recognized", zap.String("signature", ctx.Get(apps.HeaderClientSignature)),
				zap.String("code", b.Code))
			return ctx.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Meta: model.Meta{
					Code:    apps.ErrCodeUnauthorized,
					Message: apps.ErrMsgUnauthorized,
				},
			})
		}
		return ctx.Next()
	}
}
