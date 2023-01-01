package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParameter_FindByParamGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	g := "Group A"
	pool := pgxpoolmock.NewMockPgxIface(ctrl)
	ctx := context.Background()
	persister := NewParameter(Parameter{
		Logger: logger,
		Pool:   pool,
	})
	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"id", "param_group", "param_name", "param_value"}).AddRow(
			sql.NullInt64{Int64: 1, Valid: true},
			sql.NullString{String: "Group A", Valid: true},
			sql.NullString{String: "Name A", Valid: true},
			sql.NullString{String: "Value A", Valid: true},
		).ToPgxRows()
		pool.EXPECT().Query(ctx,
			"SELECT id, param_group, param_name, param_value from parameters where param_group = $1", g).Return(rows, nil)
		data, ex := persister.FindByParamGroup(g)
		assert.NotNil(t, data)
		assert.Nil(t, ex)
	})

	t.Run("should return exception on failed to execute query", func(t *testing.T) {
		pool.EXPECT().Query(ctx,
			"SELECT id, param_group, param_name, param_value from parameters where param_group = $1", g).
			Return(nil, fmt.Errorf("something went wrong on execute query"))
		data, ex := persister.FindByParamGroup(g)
		assert.Nil(t, data)
		assert.NotNil(t, ex)
	})

	t.Run("should return exception on map query result", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"id", "param_group", "param_name", "param_value"}).AddRow(
			1,
			sql.NullString{String: "Group A", Valid: true},
			sql.NullString{String: "Name A", Valid: true},
			sql.NullString{String: "Value A", Valid: true},
		).ToPgxRows()
		pool.EXPECT().Query(ctx,
			"SELECT id, param_group, param_name, param_value from parameters where param_group = $1", g).
			Return(rows, nil)
		data, ex := persister.FindByParamGroup(g)
		assert.Nil(t, data)
		assert.NotNil(t, ex)
	})

}
