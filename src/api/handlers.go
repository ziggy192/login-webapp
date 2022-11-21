package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/auth"
	"bitbucket.org/ziggy192/ng_lu/src/api/model"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"encoding/json"
	"fmt"
	"google.golang.org/api/idtoken"
	"net/http"
	"strings"
)

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var rqBody model.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&rqBody)
	if err != nil {
		_ = util.SendJSON(ctx, w, http.StatusBadRequest, "invalid request", nil)
		return
	}

	acc, err := a.DBStores.Account.FindAccountByUserName(ctx, rqBody.Username)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	if acc == nil {
		_ = util.SendJSON(ctx, w, http.StatusUnauthorized, "username entered does not exist", nil)
		return
	}

	if !acc.IsCorrectPassword(rqBody.Password) {
		_ = util.SendJSON(ctx, w, http.StatusUnauthorized, "password is incorrect", nil)
		return
	}

	signed, err := a.Authenticator.SignUserJWT(ctx, acc.Username)
	if err != nil {
		_ = util.SendError(ctx, w, err)
		return
	}
	resp := &model.TokenResponse{AccessToken: signed}
	_ = util.SendJSON(ctx, w, 200, "login successfully", resp)
}

func (a *App) handleSignup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var rqBody model.SignupRequest
	err := json.NewDecoder(r.Body).Decode(&rqBody)
	if err != nil {
		_ = util.SendJSON(r.Context(), w, http.StatusBadRequest, "invalid request", nil)
		return
	}

	if err = rqBody.Validate(); err != nil {
		_ = util.SendJSON(r.Context(), w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	acc, err := a.DBStores.Account.FindAccountByUserName(ctx, rqBody.Username)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	if acc != nil {
		msg := fmt.Sprintf("account with username %s already exists", rqBody.Username)
		logger.Info(ctx, msg)
		_ = util.SendJSON(ctx, w, http.StatusBadRequest, msg, nil)
		return
	}

	account, err := model.NewAccount(ctx, rqBody.Username, rqBody.Password)
	if err != nil {
		_ = util.SendError(r.Context(), w, err)
		return
	}

	err = a.DBStores.Account.CreateAccount(ctx, account.Username, account.HashedPassword)
	if err != nil {
		_ = util.SendError(r.Context(), w, err)
		return
	}

	jwt, err := a.Authenticator.SignUserJWT(ctx, account.Username)
	if err != nil {
		_ = util.SendError(r.Context(), w, err)
		return
	}
	resp := &model.TokenResponse{AccessToken: jwt}
	_ = util.SendJSON(ctx, w, http.StatusOK, "signed up successfully", resp)
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

	username := tokenPayload.Claims["email"].(string)
	acc, err := a.DBStores.Account.FindAccountByUserName(ctx, username)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	googleID := tokenPayload.Subject
	if acc == nil {
		logger.Info(ctx, "account", username, "not exists yet, creating new account with googleID")
		if err := a.DBStores.Account.CreateAccountGoogle(ctx, username, googleID); err != nil {
			logger.Err(ctx, err)
			_ = util.SendError(ctx, w, err)
			return
		}
	} else if len(acc.GoogleID) == 0 {
		logger.Info(ctx, "account", username, "exists, updating googleID to account")
		if err := a.DBStores.Account.UpdateGoogleIDToAccount(ctx, username, googleID); err != nil {
			logger.Err(ctx, err)
			_ = util.SendError(ctx, w, err)
			return
		}
	}

	signed, err := a.Authenticator.SignUserJWT(ctx, username)
	if err != nil {
		_ = util.SendError(ctx, w, err)
		return
	}

	resp := &model.TokenResponse{AccessToken: signed}
	_ = util.SendJSON(ctx, w, 200, "login successfully", resp)
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

func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bearerToken := r.Header.Get(AuthorizationHeader)
	tokenString := strings.TrimPrefix(bearerToken, "Bearer ")
	err := a.Authenticator.Logout(ctx, tokenString)
	if err != nil {
		_ = util.SendError(ctx, w, err)
	}
	_ = util.SendJSON(ctx, w, http.StatusOK, "logged out successfully", nil)
}
