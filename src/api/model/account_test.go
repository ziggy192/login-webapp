package model

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccount_GenerateAndCheckPassword(t *testing.T) {
	t.Parallel()
	ctx := test.NewContext(t)
	username := uuid.NewString()
	password := uuid.NewString()
	account, err := NewAccount(ctx, username, password)
	require.NoError(t, err)
	require.True(t, account.IsCorrectPassword(password))
	for j := 0; j < 10; j++ {
		require.False(t, account.IsCorrectPassword(uuid.NewString()))
	}
}
