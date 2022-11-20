package frontend

import (
	"bitbucket.org/ziggy192/ng_lu/src/frontend/config"
	"bitbucket.org/ziggy192/ng_lu/src/frontend/model"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"html/template"
	"log"
	"net/http"
	"strings"
)

const (
	loginGooglePath = "/login_google"
	profilePath     = "/profile"

	headerAuthorization  = "Authorization"
	cookieKeyAccessToken = "access_token"
)

type App struct {
	Config        *config.Config
	SchemaDecoder *schema.Decoder
	Tmpl          *template.Template
}

func NewApp() *App {
	a := &App{
		Config:        config.New(),
		SchemaDecoder: schema.NewDecoder(),
		Tmpl: template.Must(template.ParseFiles(
			"templates/login.html",
			"templates/signup.html",
			"templates/profile_edit.html",
			"templates/profile_view.html",
		)),
	}
	a.setupRoutes()
	return a
}

func (a *App) setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_ = a.Tmpl.ExecuteTemplate(w, "login.html", nil)
			return
		}

		err := r.ParseForm()
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
			return
		}
		var acc model.Account
		err = a.SchemaDecoder.Decode(&acc, r.PostForm)
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
			return
		}

		// todo do something with acc
		logger.Info(acc)

		// todo try login with account

		// returns fails if not success

		http.Redirect(w, r, "/profile/view", http.StatusFound)
	})

	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_ = a.Tmpl.ExecuteTemplate(w, "signup.html", nil)
			return
		}

		err := r.ParseForm()
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
			return
		}
		var acc model.Account
		err = a.SchemaDecoder.Decode(&acc, r.PostForm)
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
			return
		}

		// todo do something with acc
		logger.Info(acc)

		// todo try signup with account

		// returns fails if not success

		http.Redirect(w, r, "/profile/edit", http.StatusFound)
	})

	r.HandleFunc("/profile/view", func(w http.ResponseWriter, r *http.Request) {
		// todo get profile from token

		cookie, err := r.Cookie(cookieKeyAccessToken)
		if err != nil {
			logger.Info("cannot find \"access_token\" in cookie, redirect to login")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// todo serparate package -> service
		request, err := http.NewRequest("GET", a.Config.APIRoot+profilePath, nil)
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
			return
		}
		request.Header.Set(headerAuthorization, "Bearer "+cookie.Value)

		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
			return
		}

		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			logger.Info("unauthorized, redirect to login", "status", resp.StatusCode)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		if !util.StatusSuccess(resp.StatusCode) {
			logger.Err("error from server", "status code", resp.StatusCode, "body", util.ReadBody(resp.Body))
			_ = util.SendJSON(w, resp.StatusCode, "error from server", nil)
			return
		}

		var baseResponse model.BaseResponse
		err = json.NewDecoder(resp.Body).Decode(&baseResponse)
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
			return
		}

		var p model.Profile
		err = baseResponse.UnmarshalData(&p)
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
		}

		_ = a.Tmpl.ExecuteTemplate(w, "profile_view.html", p)
	}).Methods("GET")

	r.HandleFunc("/profile/edit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// todo get profile from token
			p := model.Profile{
				FullName: "Something",
				Phone:    "somephne",
				Email:    "email",
			}
			_ = a.Tmpl.ExecuteTemplate(w, "profile_edit.html", p)
			return
		}

		err := r.ParseForm()
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
			return
		}
		var p model.Profile
		err = a.SchemaDecoder.Decode(&p, r.PostForm)
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
			return
		}

		// todo do something with p
		logger.Info(p)

		// todo save p to database

		// todo returns fails if not success

		http.Redirect(w, r, "/profile/view", http.StatusFound)
		_ = a.Tmpl.ExecuteTemplate(w, "profile_edit.html", p)
	})

	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		// todo call api to logout
		http.Redirect(w, r, "/", http.StatusFound)
	})

	r.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
		}

		request, err := http.NewRequest("POST", a.Config.APIRoot+loginGooglePath, strings.NewReader(r.PostForm.Encode()))
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
			return
		}
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		csrfTokenCookie, err := r.Cookie("g_csrf_token")
		if err != nil {
			_ = util.SendJSON(w, 400, "no CSRF token in Cookie", nil)
			return
		}

		request.AddCookie(&http.Cookie{
			Name:   csrfTokenCookie.Name,
			Value:  csrfTokenCookie.Value,
			MaxAge: csrfTokenCookie.MaxAge,
		})
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
			return
		}

		if !util.StatusSuccess(resp.StatusCode) {
			logger.Err("error authentication from server", "status code", resp.StatusCode, "body", util.ReadBody(resp.Body))
			_ = util.SendJSON(w, resp.StatusCode, "error authentication from server", nil)
			return
		}

		var baseResponse model.BaseResponse
		err = json.NewDecoder(resp.Body).Decode(&baseResponse)
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  cookieKeyAccessToken,
			Value: baseResponse.Data["access_token"].(string),
		})

		http.Redirect(w, r, "/profile/view", http.StatusFound)
	})

	r.Use(util.LoggingMiddleware)
	http.Handle("/", r)
}

func (a *App) Start() {
	logger.Info("listening on port", a.Config.Port)
	log.Fatal(http.ListenAndServe(":"+a.Config.Port, nil))
}

// Stop stops app
func (a *App) Stop() {}
