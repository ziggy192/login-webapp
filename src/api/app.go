package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/config"
	"bitbucket.org/ziggy192/ng_lu/src/api/model"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"google.golang.org/api/idtoken"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	contextKeyUser     = "user"
	contextKeyIssuedAt = "issued_at"
)

type App struct {
	Config *config.Config
}

func NewApp() *App {
	a := &App{
		Config: config.New(),
	}
	a.setupRoutes()
	return a
}

func (a *App) setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

	}).Methods("POST")
	r.HandleFunc("/login_google", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("url", r.URL)
		csrfTokenCookie, err := r.Cookie("g_csrf_token")
		if err != nil {
			_ = util.SendJSON(w, 400, "no CSRF token in Cookie", nil)
			return
		}
		csrfTokenBody := r.PostFormValue("g_csrf_token")
		if len(csrfTokenBody) == 0 {
			_ = util.SendJSON(w, 400, "no CSRF token in post body", nil)
			return
		}

		if csrfTokenBody != csrfTokenCookie.Value {
			_ = util.SendJSON(w, 400, "failed to verify double submit cookie", nil)
			return
		}

		credential := r.PostFormValue("credential")
		selectBy := r.PostFormValue("select_by")
		logger.Info("credential", credential)
		logger.Info("select_by", selectBy)

		// todo should use context propagation ?
		tokenPayload, err := idtoken.Validate(context.Background(), credential, "588338350106-u3e7ddin0njjervl05577fioq678nbi5.apps.googleusercontent.com")
		if err != nil {
			logger.Err(err)
			_ = util.SendJSON(w, 401, "Invalid ID Token", nil)
			return
		}

		logger.Info("token verified", tokenPayload.Claims)

		// todo lookup database to signup or login with google account

		ac := model.Account{
			Username: tokenPayload.Claims["email"].(string),
			Password: "nghia",
			GoogleID: tokenPayload.Subject,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": ac.Username,
			"iat": time.Now().Unix(),
		})
		signedString, err := token.SignedString([]byte(a.Config.AuthSecret))
		if err != nil {
			logger.Err(err)
			_ = util.SendError(w, err)
			return
		}

		data := map[string]any{
			"access_token": signedString,
		}
		_ = util.SendJSON(w, 200, "login successfully", data)
	}).Methods("POST")

	r.HandleFunc("/signup", nil).Methods("POST")

	profileR := r.PathPrefix("/profile").Subrouter()
	profileR.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		// todo get profile from database by email

		p := &model.Profile{
			FullName: "nghia api",
			Phone:    "something",
			Email:    r.Context().Value(contextKeyUser).(string),
		}
		_ = util.SendJSON(w, http.StatusOK, "get profile successfully", p)
	}).Methods("GET") // get profile by id using the jwt token
	profileR.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		// todo get profile by email
		p := &model.Profile{
			FullName: "nghia api",
			Phone:    "something",
			Email:    r.Context().Value(contextKeyUser).(string),
		}
		_ = util.SendJSON(w, http.StatusOK, "save profile successfully", p)
	}).Methods("PUT")

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
		bearerToken := r.Header.Get("Authorization")
		if len(bearerToken) == 0 {
			_ = util.SendJSON(w, http.StatusUnauthorized, "no token found", nil)
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
			logger.Err(err)
			_ = util.SendJSON(w, http.StatusUnauthorized, "invalid token", nil)
		}

		mapClaims := token.Claims.(jwt.MapClaims)
		issuedAt := int64(mapClaims["iat"].(float64))
		// todo check if issue at after last_logout

		user := mapClaims["sub"].(string)
		logger.Info("authenticated user", user, "issued at", time.Unix(issuedAt, 0))
		r = r.WithContext(context.WithValue(r.Context(), contextKeyUser, user))
		r = r.WithContext(context.WithValue(r.Context(), contextKeyIssuedAt, issuedAt))
		next.ServeHTTP(w, r)
	})
}

func (a *App) Start() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("listening on port", a.Config.Port)
	log.Fatal(http.ListenAndServe(":"+a.Config.Port, nil))
}

// Stop stops app
func (a *App) Stop() {}
