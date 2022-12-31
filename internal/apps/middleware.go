package apps

import (
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strings"
)

type Middleware struct {
	Logger *zap.Logger
}

func unauthorized(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(model.Response{
		Meta: model.Meta{
			Code:    ErrCodeUnauthorized,
			Message: ErrMsgUnauthorized,
		},
	})
}

func Authenticator(m Middleware) fiber.Handler {
	signs := []string{
		"/api/v1/authorization/client",
	}
	return func(ctx *fiber.Ctx) error {
		isSign := false
		for _, s := range signs {
			if strings.Contains(string(ctx.Request().URI().Path()), s) {
				isSign = true
				break
			}
		}

		if isSign {
			b := struct {
				Code string `json:"code"`
			}{}
			if err := ctx.BodyParser(&b); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(err)
			}

			epoch := ctx.Get(HeaderClientTimestamp)
			key := ctx.Get(HeaderApiKey)

			d := string(ctx.Request().Header.Method()) + ":" + strings.ToUpper(b.Code) + ":" + epoch + ":" + strings.ToUpper(key)
			if HMAC(d, b.Code) != ctx.Get(HeaderClientSignature) {
				m.Logger.Error("the given signature is not recognized", zap.String("signature", ctx.Get(HeaderClientSignature)),
					zap.String("code", b.Code))
				return unauthorized(ctx)
			}
		}

		return ctx.Next()
	}
}
