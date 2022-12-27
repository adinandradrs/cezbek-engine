package storage

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"strings"
)

type PgOptions struct {
	Host    string
	Port    string
	User    string
	Passwd  string
	Schema  string
	Options *string
	Logger  *zap.Logger
}

type PgPool struct {
	Pool *pgxpool.Pool
}

type Pooler interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

func NewPgPool(pg *PgOptions) *PgPool {
	url := "postgres://{{username}}:{{password}}@{{host}}:{{port}}/{{schema}}"
	url = strings.Replace(url, "{{host}}", pg.Host, -1)
	url = strings.Replace(url, "{{port}}", pg.Port, -1)
	url = strings.Replace(url, "{{username}}", pg.User, -1)
	url = strings.Replace(url, "{{password}}", pg.Passwd, -1)
	url = strings.Replace(url, "{{schema}}", pg.Schema, -1)
	if pg.Options != nil {
		url += "?" + *pg.Options
	}
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		pg.Logger.Fatal("failed to settle postgres connection", zap.Error(err))
	}
	return &PgPool{Pool: pool}
}
