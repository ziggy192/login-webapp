package auth

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTokenBlocker_BLockAndCheck(t *testing.T) {
	t.Parallel()

	ctx := newTestCtx(t)

	b := NewTokenBlocker(redisClient)
	token := t.Name() + uuid.NewString()
	err := b.BlockToken(ctx, token, time.Now(), time.Second)
	require.NoError(t, err)

	blocked, err := b.IsBlocked(ctx, token)
	require.NoError(t, err)
	require.True(t, blocked)

	time.Sleep(time.Second + 500*time.Millisecond)
	blocked, err = b.IsBlocked(ctx, token)
	require.NoError(t, err)
	require.False(t, blocked)
}
