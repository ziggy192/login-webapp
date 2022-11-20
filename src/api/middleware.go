package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/auth"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	Authenticator *auth.Authenticator
}

func NewAuthMiddleware(a *auth.Authenticator) *AuthMiddleware {
	return &AuthMiddleware{Authenticator: a}
}

func (a *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		bearerToken := r.Header.Get("Authorization")
		if len(bearerToken) == 0 {
			_ = util.SendJSON(ctx, w, http.StatusUnauthorized, "no token found", nil)
			return
		}

		tokenString := strings.TrimPrefix(bearerToken, "Bearer ")
		claims, err := a.Authenticator.VerifyUserJWT(ctx, tokenString)
		if err != nil {
			logger.Err(ctx, err)
			_ = util.SendJSON(ctx, w, http.StatusUnauthorized, "invalid token", nil)
		}

		// todo check if issue at after last_logout
		username := claims.Subject
		logger.Info(ctx, "authenticated user", username, "claims", *claims)
		r = r.WithContext(auth.SaveClaims(r.Context(), claims))
		next.ServeHTTP(w, r)
	})
}
