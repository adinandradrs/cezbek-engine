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

func isSignedUrl(signs *[]string, ctx *fiber.Ctx) bool {
	isSign := false
	for _, s := range *signs {
		if strings.Contains(string(ctx.Request().URI().Path()), s) {
			isSign = true
			break
		}
	}
	return isSign
}

func validateNonSignedAccess(ctx *fiber.Ctx) error {
	if ctx.Get(HeaderClientChannel) != ChannelEBizKezbek {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Meta: model.Meta{
				Code:    ErrCodeInvalidChannel,
				Message: ErrMsgInvalidChannel,
			},
		})
	}
	return nil
}

func validateSignedAccess(ctx *fiber.Ctx) (*string, error) {
	if ctx.Get(HeaderClientChannel) != ChannelB2BClient {
		return nil, ctx.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Meta: model.Meta{
				Code:    ErrCodeInvalidChannel,
				Message: ErrMsgInvalidChannel,
			},
		})
	}
	b := struct {
		Code string `json:"code"`
	}{}
	if err := ctx.BodyParser(&b); err != nil {
		return nil, ctx.Status(fiber.StatusBadRequest).JSON(err)
	}
	return &b.Code, nil
}

func Authenticator(m Middleware) fiber.Handler {
	signs := []string{
		"/api/v1/authorization/client",
	}
	return func(ctx *fiber.Ctx) error {
		if isSignedUrl(&signs, ctx) {
			code, err := validateSignedAccess(ctx)
			if err != nil {
				return err
			}
			epoch := ctx.Get(HeaderClientTimestamp)
			key := ctx.Get(HeaderApiKey)

			d := string(ctx.Request().Header.Method()) + ":" + strings.ToUpper(*code) + ":" + epoch + ":" + strings.ToUpper(key)
			if HMAC(d, *code) != ctx.Get(HeaderClientSignature) {
				m.Logger.Error("the given signature is not recognized", zap.String("signature", ctx.Get(HeaderClientSignature)),
					zap.String("code", *code))
				return unauthorized(ctx)
			}
		}
		err := validateNonSignedAccess(ctx)
		if err != nil {
			return err
		}
		return ctx.Next()
	}
}
