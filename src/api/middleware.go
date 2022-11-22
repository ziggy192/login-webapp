package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/auth"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"fmt"
	"net/http"
	"strings"
)

const AuthorizationHeader = "Authorization"
const AuthorizationBearer = "Bearer"

type AuthMiddleware struct {
	Authenticator *auth.Authenticator
}

func NewAuthMiddleware(a *auth.Authenticator) *AuthMiddleware {
	return &AuthMiddleware{Authenticator: a}
}

func (a *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger.Info(ctx, "authentication middleware")
		authorizationHeader := r.Header.Get(AuthorizationHeader)
		if len(authorizationHeader) == 0 {
			_ = util.SendJSON(ctx, w, http.StatusUnauthorized, "no token found", nil)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			_ = util.SendJSON(ctx, w, http.StatusUnauthorized, "invalid authorization header format", nil)
			return
		}

		authorizationType := fields[0]
		if strings.ToLower(authorizationType) != strings.ToLower(AuthorizationBearer) {
			msg := fmt.Sprintf("unsupported authorization type %s", authorizationType)
			_ = util.SendJSON(ctx, w, http.StatusUnauthorized, msg, nil)
			return
		}

		tokenString := fields[1]
		claims, err := a.Authenticator.VerifyUserJWT(ctx, tokenString)
		if err != nil {
			logger.Err(ctx, err)
			_ = util.SendJSON(ctx, w, http.StatusUnauthorized, "invalid token", nil)
			return
		}

		username := claims.Subject
		logger.Info(ctx, "authenticated user", username, fmt.Sprintf("claims %+v", *claims))
		r = r.WithContext(auth.SaveClaims(r.Context(), claims))
		next.ServeHTTP(w, r)
	})
}
