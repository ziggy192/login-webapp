package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/auth"
	"bitbucket.org/ziggy192/ng_lu/src/api/config"
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
	Authenticator *auth.Authenticator
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := config.New()
	dbStores, err := store.NewDBStores(ctx, cfg.MySQL)
	if err != nil {
		logger.Err(ctx, err)
		return nil, err
	}

	a := &App{
		Config:        cfg,
		DBStores:      dbStores,
		Authenticator: auth.NewAuthenticator(cfg),
	}

	a.setupRoutes()
	return a, nil
}

func (a *App) setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/login", a.handleLogin).Methods(http.MethodPost)
	r.HandleFunc("/login_google", a.handleLoginGoogle).Methods(http.MethodPost)
	r.HandleFunc("/signup", a.handleSignup).Methods(http.MethodPost)

	profileR := r.PathPrefix("/profile").Subrouter()
	profileR.HandleFunc("", a.handleGetProfile).Methods(http.MethodGet)
	profileR.HandleFunc("", a.handleSaveProfile).Methods(http.MethodPut)
	profileR.Use(NewAuthMiddleware(a.Authenticator).Middleware)

	r.Use(util.LoggingMiddleware)

	http.Handle("/", r)
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
