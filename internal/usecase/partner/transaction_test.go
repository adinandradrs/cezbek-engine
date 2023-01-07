package partner

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/repository"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransaction_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	dao := repository.NewMockTransactionPersister(ctrl)
	svc := NewTransaction(Transaction{
		Dao:    dao,
		Logger: logger,
	})
	inp := &model.SearchRequest{
		Limit:  3,
		Start:  1,
		SortBy: "DESC",
		SessionRequest: model.SessionRequest{
			Id: 1,
		},
	}
	t.Run("should success", func(t *testing.T) {
		count := 50
		dao.EXPECT().SearchByPartner(inp).Return([]model.PartnerTransactionProjection{
			{
				Id:          1,
				Transaction: decimal.NewFromInt(50000),
				Qty:         1,
				WalletCode:  "WALLET_A",
				Email:       "someone1@email.net",
				Cashback:    decimal.NewFromInt(500),
				Reward:      decimal.NewFromInt(1500),
				Msisdn:      "628118770510",
			},
			{
				Id:          2,
				Transaction: decimal.NewFromInt(75000),
				Qty:         3,
				WalletCode:  "WALLET_A",
				Email:       "someone2@email.net",
				Cashback:    decimal.NewFromInt(850),
				Reward:      decimal.NewFromInt(1400),
				Msisdn:      "628118770511",
			},
			{
				Id:          3,
				Transaction: decimal.NewFromInt(45000),
				Qty:         5,
				WalletCode:  "WALLET_A",
				Email:       "someone3@email.net",
				Cashback:    decimal.NewFromInt(300),
				Reward:      decimal.NewFromInt(400),
				Msisdn:      "628118770511",
			},
		}, nil)
		dao.EXPECT().CountByPartner(inp).Return(&count, nil)
		v, ex := svc.Search(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should return exception on failed in one of the query", func(t *testing.T) {
		count := 50
		dao.EXPECT().SearchByPartner(inp).Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-001",
		})
		dao.EXPECT().CountByPartner(inp).Return(&count, nil)
		v, ex := svc.Search(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}

func TestTransaction_Detail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	dao := repository.NewMockTransactionPersister(ctrl)
	inp := &model.FindByIdRequest{
		Id: 1,
		SessionRequest: model.SessionRequest{
			Id: 1,
		},
	}
	svc := NewTransaction(Transaction{
		Dao:    dao,
		Logger: logger,
	})

	t.Run("should success", func(t *testing.T) {
		dao.EXPECT().DetailByPartner(inp).Return(&model.PartnerTransactionProjection{
			Transaction: decimal.NewFromInt(1000),
			Id:          1,
			Reward:      decimal.NewFromInt(500),
			Msisdn:      "628118770510",
			Email:       "someone@email.net",
			Cashback:    decimal.NewFromInt(350),
			Qty:         1,
			WalletCode:  "CODE_A",
		}, nil)
		v, ex := svc.Detail(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should return exception on failed to find detail", func(t *testing.T) {
		dao.EXPECT().DetailByPartner(inp).Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-001",
		})
		v, ex := svc.Detail(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}
