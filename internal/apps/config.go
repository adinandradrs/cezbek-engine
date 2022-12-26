package apps

import (
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

func NewEnv(logger *zap.Logger) (*viper.Viper, *model.TechnicalError) {
	v := viper.New()
	if _, err := os.Stat(".env"); err == nil {
		v.SetConfigFile(".env")
		err := v.ReadInConfig()
		if err != nil {
			logger.Fatal("failed to load configuration from .env file", zap.Error(err))
			return nil, &model.TechnicalError{
				Exception: err.Error(),
			}
		}
	} else {
		logger.Info("no configuration from .env")
		v.AutomaticEnv()
	}
	return v, nil
}
