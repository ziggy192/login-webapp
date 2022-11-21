package frontend

import (
	"bitbucket.org/ziggy192/ng_lu/src/frontend/config"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"context"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
)

const (
	pathLogin       = "/login"
	pathLoginGoogle = "/login_google"
	pathSignup      = "/signup"
	pathLogout      = "/logout"
	pathProfile     = "/profile"

	headerAuthorization  = "Authorization"
	cookieKeyAccessToken = "access_token"
)

type App struct {
	Config        *config.Config
	SchemaDecoder *schema.Decoder
	Tmpl          *Template
}

func NewApp() *App {
	a := &App{
		Config:        config.New(),
		SchemaDecoder: schema.NewDecoder(),
		Tmpl:          NewTemplate(),
	}
	a.setupRoutes()
	return a
}

func (a *App) setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/", a.handlePostLogin).Methods(http.MethodPost)
	r.HandleFunc("/", a.handleGetLogin).Methods(http.MethodGet)

	r.HandleFunc("/signup", a.handleGetSignup).Methods(http.MethodGet)
	r.HandleFunc("/signup", a.handlePostSignup).Methods(http.MethodPost)
	r.HandleFunc("/profile/view", a.handleProfileView).Methods(http.MethodGet)
	r.HandleFunc("/profile/edit", a.handleProfileEdit).Methods(http.MethodPost, http.MethodGet)
	r.HandleFunc("/logout", a.handleLogout).Methods(http.MethodGet)
	r.HandleFunc("/auth", a.handleAuth).Methods(http.MethodPost)

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
