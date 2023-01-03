package client

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/repository"
	"github.com/adinandradrs/cezbek-engine/mock/storage"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransaction_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	dao, cacher := repository.NewMockTransactionPersister(ctrl), storage.NewMockCacher(ctrl)
	svc := NewTransaction(Transaction{
		Logger: logger,
		Cacher: cacher,
		Dao:    dao,
	})
	inp := model.TransactionRequest{
		MerchantCode:         "WCODE_A",
		Email:                "someone@email.net",
		TransactionReference: "INV/001/002/003",
		Amount:               decimal.New(150000, 10),
		Qty:                  3,
		Msisdn:               "6281123456890",
		SessionRequest: model.SessionRequest{
			Id:          int64(1),
			Email:       "corporate@email.xyz",
			PartnerCode: "CORP_A",
			Fullname:    "PT. Corporate A",
		},
	}
	t.Run("should success", func(t *testing.T) {
		cacher.EXPECT().Hget("WALLET_CODE", inp.MerchantCode).Return("WCODE_A", nil)
		dao.EXPECT().Add(gomock.Any()).Return(nil)
		v, ex := svc.Add(&inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should error on data access failed to insert", func(t *testing.T) {
		cacher.EXPECT().Hget("WALLET_CODE", inp.MerchantCode).Return("WCODE_A", nil)
		dao.EXPECT().Add(gomock.Any()).Return(&model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-001",
		})
		v, ex := svc.Add(&inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
		assert.Equal(t, apps.ErrCodeBussClientAddTransaction, ex.ErrorCode)
	})

	t.Run("should return exception on wallet is not found", func(t *testing.T) {
		cacher.EXPECT().Hget("WALLET_CODE", inp.MerchantCode).Return("", &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-001",
		})
		v, ex := svc.Add(&inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
		assert.Equal(t, apps.ErrCodeBussMerchantCodeInvalid, ex.ErrorCode)
	})
}
