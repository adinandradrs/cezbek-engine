package client

import (
	"database/sql"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"go.uber.org/zap"
)

type Transaction struct {
	Dao    repository.TransactionPersister
	Cacher storage.Cacher
	Logger *zap.Logger
}

type TransactionProvider interface {
	Add(inp *model.TransactionRequest) (*model.TransactionResponse, *model.BusinessError)
}

func NewTransaction(t Transaction) TransactionProvider {
	return &t
}

func (t *Transaction) Add(inp *model.TransactionRequest) (*model.TransactionResponse, *model.BusinessError) {
	trx := apps.Transaction(inp.Msisdn)
	_, ex := t.Cacher.Hget("WALLET_CODE", inp.WalletCode)
	if ex != nil {
		t.Logger.Error("failed to add transaction - invalid wallet code", zap.Any("tx", inp))
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussWalletCodeInvalid,
			ErrorMessage: apps.ErrMsgBussWalletCodeInvalid,
		}
	}
	data := model.Transaction{
		PartnerId:      inp.SessionRequest.Id,
		Partner:        sql.NullString{String: inp.SessionRequest.Fullname, Valid: true},
		WalletCode:     sql.NullString{String: inp.WalletCode, Valid: true},
		Qty:            inp.Qty,
		Amount:         inp.Amount,
		Msisdn:         sql.NullString{String: inp.Msisdn, Valid: true},
		Email:          sql.NullString{String: inp.SessionRequest.Email, Valid: true},
		KezbekRefCode:  sql.NullString{String: trx.TransactionId, Valid: true},
		PartnerRefCode: sql.NullString{String: inp.TransactionReference, Valid: true},
		BaseEntity: model.BaseEntity{
			CreatedBy: sql.NullInt64{Int64: inp.SessionRequest.Id},
		},
	}
	ex = t.Dao.Add(data)
	if ex != nil {
		t.Logger.Error("failed to add transaction - data access", zap.Any("tx", inp))
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussClientAddTransaction,
			ErrorMessage: apps.ErrMsgBussClientAddTransaction,
		}
	}
	return &trx, nil
}
