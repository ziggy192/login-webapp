package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/model"
	"bitbucket.org/ziggy192/ng_lu/src/api/test"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginAPI(t *testing.T) {
	t.Parallel()
	ctx := test.NewContext(t)

	app, err := NewApp(ctx)
	require.NoError(t, err)

	acc, password := randomAccount(t)
	err = app.DBStores.Account.CreateAccount(context.Background(), acc.Username, acc.HashedPassword)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		body        map[string]any
		requireFunc func(t *testing.T, r *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]any{
				"username": acc.Username,
				"password": password,
			},
			requireFunc: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, r.Code)

				data, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				var tokenResp = &model.TokenResponse{}
				gotResponse := util.BaseResponse{Data: &tokenResp}
				err = json.Unmarshal(data, &gotResponse)
				require.NoError(t, err)
				require.NotEmpty(t, tokenResp.AccessToken)

				_, err = app.Authenticator.VerifyUserJWT(test.NewContext(t), tokenResp.AccessToken)
				require.NoError(t, err)
			},
		},
		{
			name: "UserNotFound",
			body: map[string]any{
				"username": "NotFound",
				"password": "password",
			},
			requireFunc: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, r.Code)
			},
		},
		{
			name: "IncorrectPassword",
			body: map[string]any{
				"username": acc.Username,
				"password": "incorrect",
			},
			requireFunc: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, r.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			app.router.ServeHTTP(recorder, request)
			tc.requireFunc(t, recorder)
		})
	}
}

func randomAccount(t *testing.T) (account *model.Account, password string) {
	ctx := test.NewContext(t)
	username := util.RandomString(6)
	password = util.RandomString(6)
	account, err := model.NewAccount(ctx, username, password)
	require.NoError(t, err)
	return
}
