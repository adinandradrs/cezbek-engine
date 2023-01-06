package workflow

import (
	"database/sql"
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type Tier struct {
	Dao            repository.TierPersister
	Cacher         storage.Cacher
	ExpiryDuration time.Duration
	Logger         *zap.Logger
}

type TierProvider interface {
	Save(inp *model.TierRequest) (*model.WfRewardTierProjection, *model.TechnicalError)
}

func NewTier(t Tier) TierProvider {
	return &t
}

func (t Tier) add(inp *model.TierRequest) *model.TechnicalError {
	return t.Dao.Add(model.Tier{
		PartnerId:            inp.PartnerId,
		Msisdn:               sql.NullString{String: inp.Msisdn},
		Email:                sql.NullString{String: inp.Email},
		CurrentGrade:         1,
		CurrentTier:          sql.NullString{String: "BRONZE"},
		PrevGrade:            1,
		PrevTier:             sql.NullString{String: "BRONZE"},
		TransactionRecurring: 1,
		ExpiredDate:          sql.NullTime{Time: time.Now().Add(t.ExpiryDuration)},
		Journey: model.TierJourneys{
			CurrentTier:       sql.NullString{String: "BRONZE"},
			CurrentGrade:      1,
			LastTransactionId: inp.TransactionId,
		},
		BaseEntity: model.BaseEntity{
			CreatedBy: sql.NullInt64{Int64: inp.PartnerId},
		},
	})
}

func (t Tier) update(v *model.Tier, inp *model.TierRequest) (*model.WfRewardTierProjection, *model.TechnicalError) {
	v.TransactionRecurring = v.TransactionRecurring + 1
	cacher, ex := t.Cacher.Get("WFREWARD:"+v.CurrentTier.String, strconv.Itoa(v.TransactionRecurring))
	var m model.WfRewardTierProjection
	currentRecurring := v.TransactionRecurring
	if ex == nil && cacher != "" {
		_ = json.Unmarshal([]byte(cacher), &m)
		t.Logger.Info("workflow tier found", zap.String("msisdn", inp.Msisdn), zap.Any("reward", m.Reward))
		if m.MaxRecurring == v.TransactionRecurring && m.NextTier.Grade != nil {
			v.PrevGrade = v.CurrentGrade
			v.PrevTier = v.CurrentTier
			v.CurrentGrade = *m.NextTier.Grade
			v.CurrentTier = sql.NullString{String: *m.NextTier.Tier}
			v.TransactionRecurring = 1
			v.ExpiredDate = sql.NullTime{Time: time.Now().Add(t.ExpiryDuration)}
		}
	}
	v.Journey = model.TierJourneys{
		CurrentTier:       v.CurrentTier,
		CurrentGrade:      v.CurrentGrade,
		LastTransactionId: inp.TransactionId,
	}
	v.BaseEntity.UpdatedBy = sql.NullInt64{Int64: inp.PartnerId}
	ex = t.Dao.Update(*v)
	if m.Recurring == currentRecurring ||
		m.MaxRecurring == currentRecurring {
		return &m, nil
	}
	return nil, ex
}

func (t Tier) Save(inp *model.TierRequest) (*model.WfRewardTierProjection, *model.TechnicalError) {
	v, ex := t.Dao.FindByPartnerMsisdn(inp.PartnerId, inp.Msisdn)
	if v == nil && ex != nil {
		ex = t.add(inp)
		return nil, ex
	} else {
		return t.update(v, inp)
	}
}
