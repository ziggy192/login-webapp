package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/auth"
	"bitbucket.org/ziggy192/ng_lu/src/api/config"
	"bitbucket.org/ziggy192/ng_lu/src/api/redis"
	"bitbucket.org/ziggy192/ng_lu/src/api/store"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"context"
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	Config        *config.Config
	DBStores      *store.DBStores
	RedisClient   *redis.Redis
	Authenticator *auth.Authenticator
	router        *mux.Router
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := config.New()

	redisClient, err := redis.CreateRedisClient(ctx, cfg.Redis)
	if err != nil {
		logger.Err(ctx, err)
		return nil, err
	}

	dbStores, err := store.NewDBStores(ctx, cfg.MySQL)
	if err != nil {
		logger.Err(ctx, err)
		return nil, err
	}

	a := &App{
		Config:        cfg,
		DBStores:      dbStores,
		RedisClient:   redisClient,
		Authenticator: auth.NewAuthenticator(cfg, redisClient),
	}

	a.setupRoutes()
	return a, nil
}

func (a *App) setupRoutes() {
	authMiddleware := NewAuthMiddleware(a.Authenticator)

	r := mux.NewRouter()
	r.HandleFunc("/ping", a.handlePing).Methods(http.MethodGet)
	r.HandleFunc("/login", a.handleLogin).Methods(http.MethodPost)
	r.HandleFunc("/login_google", a.handleLoginGoogle).Methods(http.MethodPost)
	r.HandleFunc("/signup", a.handleSignup).Methods(http.MethodPost)
	r.Handle("/logout", authMiddleware.Middleware(http.HandlerFunc(a.handleLogout))).Methods(http.MethodPost)

	profileR := r.PathPrefix("/profile").Subrouter()
	profileR.HandleFunc("", a.handleGetProfile).Methods(http.MethodGet)
	profileR.HandleFunc("", a.handleSaveProfile).Methods(http.MethodPut)
	profileR.Use(authMiddleware.Middleware)

	r.Use(util.RequestIDMiddleware)
	r.Use(util.LoggingMiddleware)

	a.router = r
	//http.Handle("/", r)
}

func (a *App) Start(ctx context.Context) error {
	err := a.DBStores.Ping(ctx)
	if err != nil {
		logger.Err(ctx, err)
		return err
	}

	logger.Info(context.Background(), "listening on port", a.Config.Port)
	return http.ListenAndServe(":"+a.Config.Port, a.router)
}

// Stop stops app
func (a *App) Stop(ctx context.Context) {
	if a.RedisClient != nil {
		if err := a.RedisClient.Disconnect(ctx); err != nil {
			logger.Err(ctx, err)
		}
	}

	if a.DBStores != nil {
		if err := a.DBStores.Close(); err != nil {
			logger.Err(ctx, err)
		}
	}
}

func (a *App) handlePing(writer http.ResponseWriter, request *http.Request) {
	_ = util.SendJSON(request.Context(), writer, http.StatusOK, "pong", map[string]string{
		"host": request.Host,
		"port": a.Config.Port,
	})
}
