package storage

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"time"
)

type (
	RedisOptions struct {
		Addr   string
		Addrs  []string
		Passwd string
		Index  int
		Pool   int
		Idle   int
		Logger *zap.Logger
	}

	singleRedis struct {
		logger *zap.Logger
		cache  *redis.Client
	}

	clusterRedis struct {
		logger *zap.Logger
		cache  *redis.ClusterClient
	}
)

type Cacher interface {
	Set(k string, p string, v interface{}, d time.Duration) *model.TechnicalError
	Hset(k string, p string, v interface{}) *model.TechnicalError
	Delete(k string, p string) *model.TechnicalError
	Get(k string, p string) (v string, e *model.TechnicalError)
	Hget(k string, p string) (v string, e *model.TechnicalError)
	Ttl(k string, p string) (t time.Duration, e *model.TechnicalError)
}

func NewRedis(o *RedisOptions) Cacher {
	return &singleRedis{
		logger: o.Logger,
		cache: redis.NewClient(&redis.Options{
			Addr:         o.Addr,
			Password:     o.Passwd,
			DB:           o.Index,
			PoolSize:     o.Pool,
			MinIdleConns: o.Idle,
		}),
	}
}

func NewClusterRedis(o *RedisOptions) Cacher {
	return &clusterRedis{
		logger: o.Logger,
		cache: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        o.Addrs,
			Password:     o.Passwd,
			PoolSize:     o.Pool,
			MinIdleConns: o.Idle,
		}),
	}
}

func (r *clusterRedis) Set(k string, p string, v interface{}, d time.Duration) *model.TechnicalError {
	r.cache.Del(k + ":" + p)
	if d != 0*time.Second {
		_, err := r.cache.SetNX(k+":"+p, v, d).Result()
		if err != nil {
			return apps.Exception("failed on cluster setnx ops", err, zap.String("keypair", k+":"+p), r.logger)
		}
	} else {
		_, err := r.cache.Set(k+":"+p, v, 0).Result()
		if err != nil {
			return apps.Exception("failed on cluster set ops", err, zap.String("keypair", k+":"+p), r.logger)
		}
	}
	return nil
}

func (r *clusterRedis) Hset(k string, p string, v interface{}) *model.TechnicalError {
	r.cache.HDel(k + ":" + p)
	_, err := r.cache.HSet(k, p, v).Result()
	if err != nil {
		return apps.Exception("failed on cluster hset ops", err, zap.String("keypair", k+":"+p), r.logger)
	}

	return nil
}

func (r *clusterRedis) Delete(k string, p string) (out *model.TechnicalError) {
	if cmd := r.cache.Del(k + ":" + p); cmd.Err() != nil {
		return apps.Exception("failed on cluster delete ops", cmd.Err(), zap.String("keypair", k+":"+p), r.logger)
	}
	return nil
}

func (r *clusterRedis) Get(k string, p string) (v string, e *model.TechnicalError) {
	v, err := r.cache.Get(k + ":" + p).Result()
	if err != nil {
		return v, apps.Exception("failed on cluster get ops", err, zap.String("keypair", k+":"+p), r.logger)
	}
	return v, nil
}

func (r *clusterRedis) Hget(k string, p string) (v string, e *model.TechnicalError) {
	v, err := r.cache.HGet(k, p).Result()
	if err != nil {
		return v, apps.Exception("failed on cluster hget ops", err, zap.String("keypair", k+":"+p), r.logger)
	}
	return v, nil
}

func (r *clusterRedis) Ttl(k string, p string) (t time.Duration, e *model.TechnicalError) {
	if cmd := r.cache.TTL(k + ":" + p); cmd.Err() != nil {
		r.logger.Error("failed on cluster TTL ops property", zap.String("key", k), zap.String("pair", p))
		return t, apps.Exception("failed on cluster TTL ops", cmd.Err(), zap.String("keypair", k+":"+p), r.logger)
	} else {
		return cmd.Val(), nil
	}
}

func (r *singleRedis) Set(k string, p string, v interface{}, d time.Duration) *model.TechnicalError {
	r.cache.Del(k + ":" + p)
	if d != 0*time.Second {
		_, err := r.cache.SetNX(k+":"+p, v, d).Result()
		if err != nil {
			return apps.Exception("failed on setnx ops", err, zap.String("keypair", k+":"+p), r.logger)
		}
	}
	_, err := r.cache.Set(k+":"+p, v, 0).Result()
	if err != nil {
		return apps.Exception("failed on set ops", err, zap.String("keypair", k+":"+p), r.logger)
	}
	return nil
}

func (r *singleRedis) Delete(k string, p string) *model.TechnicalError {
	if cmd := r.cache.Del(k + ":" + p); cmd.Err() != nil {
		return apps.Exception("failed on delete ops", cmd.Err(), zap.String("keypair", k+":"+p), r.logger)
	}
	return nil
}

func (r *singleRedis) Get(k string, p string) (string, *model.TechnicalError) {
	v, err := r.cache.Get(k + ":" + p).Result()
	if err != nil {
		return v, apps.Exception("failed on get ops", err, zap.String("keypair", k+":"+p), r.logger)
	}
	return v, nil
}

func (r *singleRedis) Ttl(k string, p string) (t time.Duration, e *model.TechnicalError) {
	cmd := r.cache.TTL(k + ":" + p)
	if cmd.Err() != nil {
		return t, apps.Exception("failed on TTL ops", cmd.Err(), zap.String("keypair", k+":"+p), r.logger)
	} else {
		return cmd.Val(), nil
	}
}

func (r *singleRedis) Hset(k string, p string, v interface{}) *model.TechnicalError {
	r.cache.HDel(k + ":" + p)
	_, err := r.cache.HSet(k, p, v).Result()
	if err != nil {
		return apps.Exception("failed on single hset ops", err, zap.String("keypair", k+":"+p), r.logger)
	}

	return nil
}

func (r *singleRedis) Hget(k string, p string) (v string, e *model.TechnicalError) {
	v, err := r.cache.HGet(k, p).Result()
	if err != nil {
		return v, apps.Exception("failed on single hget ops", err, zap.String("keypair", k+":"+p), r.logger)
	}
	return v, nil
}
