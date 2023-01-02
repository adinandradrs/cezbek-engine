package main

import (
	"github.com/adinandradrs/cezbek-engine/internal"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

func main() {
	c := internal.NewContainer("app_cezbek_api")
	epoch := time.Now().Unix()

	d := fiber.MethodPost + ":" + strings.ToUpper("LAJADA") + ":" + strconv.FormatInt(epoch, 10) + ":" +
		strings.ToUpper("ee33c45e2cfe3e08d352698d31da6bee")
	c.Logger.Info("EPOCH", zap.Int64("unixts", epoch), zap.String("hmac", apps.HMAC(d, "LAJADA")))
}
