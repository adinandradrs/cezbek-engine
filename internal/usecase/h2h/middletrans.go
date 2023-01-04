package h2h

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"strings"
	"time"
)

func NewMiddletrans(middletrans Middletrans) FactoryProvider {
	return &middletrans
}

func (m *Middletrans) SendCashback(inp *model.H2HSendCashbackRequest) (*model.TransactionResponse, *model.BusinessError) {
	inp.WalletCode = strings.ToLower(inp.WalletCode)
	v, ex := m.WalletTransfer(&model.MiddletransWalletTransferRequest{
		Amount:  inp.Amount,
		Account: inp.Destination,
		Wallet:  inp.WalletCode,
	})
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussH2HCashbackFailed,
			ErrorMessage: apps.ErrMsgBussH2HCashbackFailed,
		}
	}
	return &model.TransactionResponse{
		TransactionId:        v.TransactionRef,
		TransactionTimestamp: time.Now().Unix(),
	}, nil
}
