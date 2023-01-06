package job

import (
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"go.uber.org/zap"
	"time"
)

type Tier struct {
	Dao     repository.TierPersister
	Expired *time.Duration
	Logger  *zap.Logger
}

type TierWatcher interface {
	Expire()
}

func NewTier(t Tier) TierWatcher {
	return &t
}

func (t *Tier) Expire() {
	count, ex := t.Dao.CountExpire()
	if ex != nil {
		t.Logger.Error("failed to count tier that need to be expired")
	} else if ex == nil && *count > 0 {
		exp := time.Now().Add(*t.Expired)
		t.Logger.Info("expired tiers total data", zap.Int("total", *count), zap.Time("next", exp))
		ex = t.Dao.Expire(exp)
		if ex != nil {
			t.Logger.Error("failed to expire tier data")
		}
	} else {
		t.Logger.Info("expired tiers total data", zap.Int("total", *count))
	}
}
