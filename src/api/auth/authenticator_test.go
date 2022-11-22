package auth

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/config"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAuthenticator_SignAndVerifyJWT(t *testing.T) {
	t.Parallel()

	ctx := newTestCtx(t)

	validDuration := 60 * time.Minute

	const noneToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJuZ19sdSIsInN1YiI6InVzZXJuYW1lX2YzMzQyNzBiLWFiZjgtNDM2Yi1iMGIyLTRmNjJjYmMxZGQyMiIsImV4cCI6MTY2MDAwNTA0MCwiaWF0IjoxNjYwMDAzMDQwLCJqdGkiOiIwYjE1NzFjMC1jMzZjLTQxMTEtYjEwZi04OGQ2MDMwYTM2ZTMifQ."
	const expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJuZ19sdSIsInN1YiI6InVzZXJuYW1lX2YzMzQyNzBiLWFiZjgtNDM2Yi1iMGIyLTRmNjJjYmMxZGQyMiIsImV4cCI6MTY2MDAwNTA0MCwiaWF0IjoxNjYwMDAzMDQwLCJqdGkiOiIwYjE1NzFjMC1jMzZjLTQxMTEtYjEwZi04OGQ2MDMwYTM2ZTMifQ.kQxjrnmLO5y_0wx2xWCrfMKQeQ2UbayvNm9nDzQ0bIw"
	const invalidSignatureToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJuZ19sdSIsInN1YiI6InVzZXJuYW1lX2YzMzQyNzBiLWFiZjgtNDM2Yi1iMGIyLTRmNjJjYmMxZGQyMiIsImV4cCI6MTY2MDAwNTA0MCwiaWF0IjoxNjYwMDAzMDQwLCJqdGkiOiIwYjE1NzFjMC1jMzZjLTQxMTEtYjEwZi04OGQ2MDMwYTM2ZTMifQ.kQxjrnmLO5y_0wx2xWCrfKKQeQ2UbayvNm9nDzQ0bIw"

	authenticator := NewAuthenticator(config.New(), redisClient)
	authenticator.expiresAfterMinutes = 60
	username := "username_" + uuid.NewString()
	validTokenString, err := authenticator.SignUserJWT(ctx, username)
	require.NoError(t, err)

	testCases := []struct {
		Name        string
		TokenString string
		ExpectValid bool
	}{
		{"valid", validTokenString, true},
		{"invalid", invalidSignatureToken, false},
		{"expired", expiredToken, false},
		{"none_algorithm", noneToken, false},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			claims, err := authenticator.VerifyUserJWT(ctx, tc.TokenString)
			require.NotNil(t, claims)

			if tc.ExpectValid {
				require.NoError(t, err)
				require.NoError(t, claims.Valid())
				expectedExpiresAt := time.Now().Add(validDuration).Unix()
				require.InDelta(t, expectedExpiresAt, claims.ExpiresAt.Unix(), 5)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

func TestAuthenticator_Logout(t *testing.T) {
	t.Parallel()
	ctx := newTestCtx(t)

	authenticator := NewAuthenticator(config.New(), redisClient)
	authenticator.expiresAfterMinutes = 60
	username := "username_" + uuid.NewString()
	tokenString, err := authenticator.SignUserJWT(ctx, username)
	require.NoError(t, err)

	_, err = authenticator.VerifyUserJWT(ctx, tokenString)
	require.NoError(t, err)

	err = authenticator.Logout(ctx, tokenString)
	require.NoError(t, err)

	_, err = authenticator.VerifyUserJWT(ctx, tokenString)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrTokenLoggedOut)
}

func newTestCtx(t *testing.T) context.Context {
	return logger.SaveRequestID(context.Background(), t.Name())
}
