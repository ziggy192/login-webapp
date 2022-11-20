package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/auth"
	"bitbucket.org/ziggy192/ng_lu/src/api/model"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"encoding/json"
	"google.golang.org/api/idtoken"
	"net/http"
)

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var loginR model.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginR)
	if err != nil {
		_ = util.SendJSON(ctx, w, http.StatusBadRequest, "invalid request", nil)
		return
	}

	acc, err := a.DBStores.Account.FindAccountByUserName(ctx, loginR.Username)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	if acc == nil {
		_ = util.SendJSON(ctx, w, http.StatusUnauthorized, "username entered does not exist", nil)
		return
	}

	if acc.IsCorrectPassword(loginR.Password) {
		_ = util.SendJSON(ctx, w, http.StatusUnauthorized, "password is incorrect", nil)
		return
	}

}

func (a *App) handleSignup(w http.ResponseWriter, r *http.Request) {
	var signupR model.SignupRequest
	err := json.NewDecoder(r.Body).Decode(&signupR)
	if err != nil {
		_ = util.SendJSON(r.Context(), w, http.StatusBadRequest, "invalid request", nil)
		return
	}

	if err = signupR.Validate(); err != nil {
		_ = util.SendJSON(r.Context(), w, http.StatusBadRequest, err.Error(), nil)
		return
	}
}

func (a *App) handleLoginGoogle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info(ctx, "url", r.URL)
	csrfTokenCookie, err := r.Cookie("g_csrf_token")
	if err != nil {
		_ = util.SendJSON(ctx, w, 400, "no CSRF token in Cookie", nil)
		return
	}
	csrfTokenBody := r.PostFormValue("g_csrf_token")
	if len(csrfTokenBody) == 0 {
		_ = util.SendJSON(ctx, w, 400, "no CSRF token in post body", nil)
		return
	}

	if csrfTokenBody != csrfTokenCookie.Value {
		_ = util.SendJSON(ctx, w, 400, "failed to verify double submit cookie", nil)
		return
	}

	credential := r.PostFormValue("credential")
	selectBy := r.PostFormValue("select_by")
	logger.Info(ctx, "credential", credential)
	logger.Info(ctx, "select_by", selectBy)

	tokenPayload, err := idtoken.Validate(ctx, credential, "588338350106-u3e7ddin0njjervl05577fioq678nbi5.apps.googleusercontent.com")
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendJSON(ctx, w, 401, "Invalid ID Token", nil)
		return
	}

	logger.Info(ctx, "token verified", tokenPayload.Claims)

	// todo lookup database to signup or login with google account
	email := tokenPayload.Claims["email"].(string)

	acc, err := a.DBStores.Account.FindAccountByUserName(ctx, email)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	if acc == nil {
		googleID := tokenPayload.Subject
		if err := a.DBStores.Account.CreateAccountGoogle(ctx, email, googleID); err != nil {
			logger.Err(ctx, err)
			_ = util.SendError(ctx, w, err)
			return
		}

		acc, err = a.DBStores.Account.FindAccountByUserName(ctx, email)
		if err != nil {
			logger.Err(ctx, err)
			_ = util.SendError(ctx, w, err)
			return
		}
	}

	signed, err := a.Authenticator.SignUserJWT(ctx, acc.Username)
	if err != nil {
		_ = util.SendError(ctx, w, err)
		return
	}

	data := map[string]any{
		"access_token": signed,
	}
	_ = util.SendJSON(ctx, w, 200, "login successfully", data)
}

func (a *App) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	username := auth.GetUsername(r.Context())
	// todo get profile from database by username
	p := &model.Profile{
		FullName: "nghia api",
		Phone:    "something",
		Email:    username,
	}
	_ = util.SendJSON(r.Context(), w, http.StatusOK, "get profile successfully", p)
}

func (a *App) handleSaveProfile(w http.ResponseWriter, r *http.Request) {
	username := auth.GetUsername(r.Context())
	// todo get profile by username
	p := &model.Profile{
		FullName: "nghia api",
		Phone:    "something",
		Email:    username,
	}
	_ = util.SendJSON(r.Context(), w, http.StatusOK, "save profile successfully", p)
}
