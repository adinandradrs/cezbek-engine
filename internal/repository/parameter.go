package repository

import (
	"context"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/georgysavva/scany/pgxscan"
	"go.uber.org/zap"
)

type Parameter struct {
	Pool   storage.Pooler
	Logger *zap.Logger
}

type ParamPersister interface {
	FindByParamGroup(g string) ([]*model.Parameter, *model.TechnicalError)
}

func NewParameter(p Parameter) ParamPersister {
	return &p
}

func (p Parameter) FindByParamGroup(g string) ([]*model.Parameter, *model.TechnicalError) {
	var params []*model.Parameter
	rows, err := p.Pool.Query(context.Background(),
		"SELECT id, param_group, param_name, param_value from parameters where param_group = $1", g)
	if err != nil {
		return nil, apps.Exception("failed to find by param group", err, zap.String("group", g), p.Logger)
	}
	err = pgxscan.ScanAll(&params, rows)
	if err != nil {
		return nil, apps.Exception("failed to map find by param group", err, zap.String("group", g), p.Logger)
	}
	return params, nil
}
