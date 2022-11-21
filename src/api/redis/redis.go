package redis

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/config"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"context"
	"crypto/tls"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	redisConnectTimeout = 5 * time.Second
	redisReadTimeout    = 5 * time.Second
	redisWriteTimeout   = 5 * time.Second
	redisIdleTimeout    = 5 * time.Minute
)

// Redis stores a client to connect to redis
type Redis struct {
	Client *redis.Client
}

// CreateRedisClient creates redis client to connect to redis
func CreateRedisClient(ctx context.Context, config *config.RedisConfig) (*Redis, error) {
	if config.PoolSize <= 0 {
		err := errors.New("redis pool size should be positive")
		return nil, err
	}

	r := Redis{}
	cfg := &redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		DialTimeout:  redisConnectTimeout,
		ReadTimeout:  redisReadTimeout,
		WriteTimeout: redisWriteTimeout,
		IdleTimeout:  redisIdleTimeout,
	}

	if config.EnableTLS {
		cfg.TLSConfig = &tls.Config{}
	}
	r.Client = redis.NewClient(cfg).WithContext(ctx)

	_, err := r.Client.Ping(ctx).Result()
	if err != nil {
		logger.Err(ctx, "cannot connect to redis at", config.Addr, err)
		return nil, err
	}

	logger.Info(ctx, "connected to redis at", config.Addr)
	return &r, nil
}

// Disconnect closes redis client and releases any open resources.
func (r *Redis) Disconnect(ctx context.Context) error {
	err := r.Client.Close()
	if err != nil {
		logger.Err(ctx, err)
	} else {
		logger.Info(ctx, "closed connection to redis")
	}

	return err
}
