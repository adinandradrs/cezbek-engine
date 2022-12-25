package apps

import (
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

func NewEnv(logger *zap.Logger) *model.TechnicalError {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			logger.Error("failed to load configuration from .env file", zap.Error(err))
			return &model.TechnicalError{
				Exception: err.Error(),
			}
		}
	} else {
		logger.Info("no configuration from .env")
	}
	return nil
}
