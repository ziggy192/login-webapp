package frontend

import (
	"bitbucket.org/ziggy192/ng_lu/src/frontend/api"
	"bitbucket.org/ziggy192/ng_lu/src/frontend/config"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"context"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	cookieKeyAccessToken = "access_token"
)

type App struct {
	Config    *config.Config
	Tmpl      *Template
	APIClient *api.Client
}

func NewApp() *App {
	cfg := config.New()
	a := &App{
		Config:    cfg,
		Tmpl:      NewTemplate(),
		APIClient: api.NewClient(cfg.APIRoot),
	}
	a.setupRoutes()
	return a
}

func (a *App) setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/ping", a.handlePing).Methods(http.MethodGet)
	r.HandleFunc("/", a.handlePostLogin).Methods(http.MethodPost)
	r.HandleFunc("/", a.handleGetLogin).Methods(http.MethodGet)
	r.HandleFunc("/signup", a.handleGetSignup).Methods(http.MethodGet)
	r.HandleFunc("/signup", a.handlePostSignup).Methods(http.MethodPost)
	r.HandleFunc("/logout", a.handleLogout).Methods(http.MethodGet)
	r.HandleFunc("/auth", a.handleAuth).Methods(http.MethodPost)

	r.HandleFunc("/profile/view", a.handleProfileView).Methods(http.MethodGet)
	r.HandleFunc("/profile/edit", a.handleGetProfileEdit).Methods(http.MethodGet)
	r.HandleFunc("/profile/edit", a.handlePostProfileEdit).Methods(http.MethodPost)

	r.Use(util.RequestIDMiddleware)
	r.Use(util.LoggingMiddleware)
	http.Handle("/", r)
}

func (a *App) Start(ctx context.Context) error {
	logger.Info(ctx, "listening on port", a.Config.Port)
	return http.ListenAndServe(":"+a.Config.Port, nil)
}

// Stop stops app
func (a *App) Stop() {}

func (a *App) handlePing(writer http.ResponseWriter, request *http.Request) {
	_ = util.SendJSON(request.Context(), writer, http.StatusOK, "pong", map[string]string{
		"host": request.Host,
		"port": a.Config.Port,
	})
}
