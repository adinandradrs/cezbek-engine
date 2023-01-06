package workflow

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type Cashback struct {
	Dao    repository.WorkflowPersister
	Logger *zap.Logger
}

type CashbackProvider interface {
	FindCashbackAmount(inp *model.FindCashbackRequest) (*model.FindCashbackResponse, *model.BusinessError)
}

func NewCashback(c Cashback) CashbackProvider {
	return &c
}

func (c *Cashback) FindCashbackAmount(inp *model.FindCashbackRequest) (*model.FindCashbackResponse, *model.BusinessError) {
	d, ex := c.Dao.FindCashbackByTransaction(inp.Qty, inp.Amount)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussNoCashback,
			ErrorMessage: apps.ErrMsgBussNoCashback,
		}
	}
	t := inp.Amount.Mul(*d).Div(decimal.NewFromInt(100))
	return &model.FindCashbackResponse{
		Amount: t,
	}, nil
}
