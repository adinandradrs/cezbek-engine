package management

import (
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"go.uber.org/zap"
	"strconv"
)

type Workflow struct {
	Dao repository.WorkflowPersister
	storage.Cacher
	Logger *zap.Logger
}

type WorkflowManager interface {
	CacheRewardTiers() *model.TechnicalError
}

func NewWorkflow(w Workflow) WorkflowManager {
	return &w
}

func (w *Workflow) CacheRewardTiers() *model.TechnicalError {
	v, ex := w.Dao.FindRewardTiers()
	if ex != nil {
		return ex
	}
	for i := range v {
		cache, _ := json.Marshal(v[i])
		_ = w.Cacher.Set("WFREWARD:"+v[i].Tier, strconv.Itoa(v[i].Recurring), cache, 0)
	}
	return nil
}
