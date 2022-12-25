package apps

import (
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

func Exception(msg string, err error, d zap.Field, logger *zap.Logger) *model.TechnicalError {
	e := &model.TechnicalError{
		Exception: err.Error(),
		Occurred:  time.Now().Unix(),
		Ticket:    uuid.New().String(),
	}
	logger.Error(msg, zap.Any("", e), d)
	return e
}

func TransactionId(identifier string) string {
	return "C002" + time.Now().Format("060102150405") + identifier[0:5] + "0"
}

func StringExists(key string, strs []string) bool {
	for _, v := range strs {
		if v == key {
			return true
		}
	}
	return false
}

func ValidateStruct(err error) []*model.BadPayloadResponse {
	var errors []*model.BadPayloadResponse
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element model.BadPayloadResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
