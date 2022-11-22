package auth

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/config"
	"bitbucket.org/ziggy192/ng_lu/src/api/redis"
	"context"
	"os"
	"testing"
)

var redisClient *redis.Redis

func TestMain(m *testing.M) {
	var err error
	ctx := context.Background()
	redisClient, err = redis.CreateRedisClient(ctx, config.New().Redis)
	if err != nil {
		panic(err)
	}

	if err := redisClient.Client.FlushAll(ctx).Err(); err != nil {
		panic(err)
	}

	code := m.Run()
	_ = redisClient.Disconnect(ctx)
	os.Exit(code)
}
