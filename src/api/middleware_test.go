package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/auth"
	"bitbucket.org/ziggy192/ng_lu/src/api/test"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, authenticator *auth.Authenticator)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			setupAuth: func(t *testing.T, request *http.Request, authenticator *auth.Authenticator) {
				ctx := test.NewContext(t)
				token, err := authenticator.SignUserJWT(ctx, "user")
				require.NoError(t, err)

				authorizationHeader := fmt.Sprintf("%s %s", "Bearer", token)
				request.Header.Set(AuthorizationHeader, authorizationHeader)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, authenticator *auth.Authenticator) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, authenticator *auth.Authenticator) {
				ctx := test.NewContext(t)
				token, err := authenticator.SignUserJWT(ctx, "user")
				require.NoError(t, err)

				authorizationHeader := fmt.Sprintf("%s %s", "unsupported", token)
				request.Header.Set(AuthorizationHeader, authorizationHeader)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, authenticator *auth.Authenticator) {
				ctx := test.NewContext(t)
				token, err := authenticator.SignUserJWT(ctx, "user")
				require.NoError(t, err)
				request.Header.Set(AuthorizationHeader, token)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidToken",
			setupAuth: func(t *testing.T, request *http.Request, authenticator *auth.Authenticator) {
				authorizationHeader := fmt.Sprintf("%s %s", "Bearer", "invalid token")
				request.Header.Set(AuthorizationHeader, authorizationHeader)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "LogoutToken",
			setupAuth: func(t *testing.T, request *http.Request, authenticator *auth.Authenticator) {
				ctx := test.NewContext(t)

				token, err := authenticator.SignUserJWT(ctx, "user")
				require.NoError(t, err)
				err = authenticator.Logout(ctx, token)
				require.NoError(t, err)

				authorizationHeader := fmt.Sprintf("%s %s", "Bearer", token)
				request.Header.Set(AuthorizationHeader, authorizationHeader)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx := test.NewContext(t)
			app, err := NewApp(ctx)
			require.NoError(t, err)
			authMiddleware := NewAuthMiddleware(app.Authenticator)
			authPath := "/test_auth"

			app.router.Handle(authPath,
				authMiddleware.Middleware(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					_ = util.SendJSON(ctx, writer, http.StatusOK, "success", nil)
				}))).Methods(http.MethodGet)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, app.Authenticator)
			app.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
