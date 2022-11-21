package auth

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/redis"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"context"
	"errors"
	goredis "github.com/go-redis/redis/v8"
	"time"
)

const namespace = "ng_lu:token_blocker"
const timeFormat = time.RFC3339

type TokenBlocker struct {
	redis *redis.Redis
}

func NewTokenBlocker(redis *redis.Redis) *TokenBlocker {
	return &TokenBlocker{redis: redis}
}

func (b *TokenBlocker) BlockToken(ctx context.Context, token string, at time.Time, duration time.Duration) error {
	key := util.BuildKey(namespace, token)
	err := b.redis.Client.Set(ctx, key, at.Format(timeFormat), duration).Err()
	if err != nil {
		logger.Err(ctx, err)
		return err
	}
	return nil
}

func (b *TokenBlocker) IsBlocked(ctx context.Context, token string) (bool, error) {
	err := b.redis.Client.Get(ctx, util.BuildKey(namespace, token)).Err()
	if errors.Is(err, goredis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
