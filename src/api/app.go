package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/config"
	"bitbucket.org/ziggy192/ng_lu/src/api/store"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

const (
	contextKeyUser     = "user"
	contextKeyIssuedAt = "issued_at"
)

type App struct {
	Config   *config.Config
	DBStores *store.DBStores
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := config.New()
	dbStores, err := store.NewDBStores(ctx, cfg.MySQL)
	if err != nil {
		logger.Err(ctx, err)
		return nil, err
	}

	a := &App{
		Config:   cfg,
		DBStores: dbStores,
	}

	a.setupRoutes()
	return a, nil
}

func (a *App) setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/login", a.handleLogin).Methods("POST")
	r.HandleFunc("/login_google", a.handleLoginGoogle).Methods("POST")
	r.HandleFunc("/signup", a.handleSignup).Methods("POST")

	profileR := r.PathPrefix("/profile").Subrouter()
	profileR.HandleFunc("", a.handleGetProfile).Methods("GET") // get profile by id using the jwt token
	profileR.HandleFunc("", a.handleSaveProfile).Methods("PUT")
	middleware := &AuthMiddleware{Secret: []byte(a.Config.AuthSecret)}
	profileR.Use(middleware.Middleware)

	r.Use(util.LoggingMiddleware)

	http.Handle("/", r)
}

type AuthMiddleware struct {
	Secret []byte
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
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return a.Secret, nil
		})

		if err != nil || !token.Valid {
			logger.Err(ctx, err)
			_ = util.SendJSON(ctx, w, http.StatusUnauthorized, "invalid token", nil)
		}

		mapClaims := token.Claims.(jwt.MapClaims)
		issuedAt := int64(mapClaims["iat"].(float64))
		// todo check if issue at after last_logout

		user := mapClaims["sub"].(string)
		logger.Info(ctx, "authenticated user", user, "issued at", time.Unix(issuedAt, 0))
		r = r.WithContext(context.WithValue(r.Context(), contextKeyUser, user))
		r = r.WithContext(context.WithValue(r.Context(), contextKeyIssuedAt, issuedAt))
		next.ServeHTTP(w, r)
	})
}

func (a *App) Start(ctx context.Context) error {
	err := a.DBStores.Ping(ctx)
	if err != nil {
		logger.Err(ctx, err)
		return err
	}

	logger.Info(context.Background(), "listening on port", a.Config.Port)
	return http.ListenAndServe(":"+a.Config.Port, nil)
}

// Stop stops app
func (a *App) Stop(ctx context.Context) {
	if a.DBStores != nil {
		if err := a.DBStores.Close(); err != nil {
			logger.Err(ctx, err)
		}
	}
}
