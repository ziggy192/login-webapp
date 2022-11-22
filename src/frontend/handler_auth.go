package frontend

import (
	"bitbucket.org/ziggy192/ng_lu/src/frontend/model"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"fmt"
	"io"
	"net/http"
)

func (a *App) handleGetLogin(w http.ResponseWriter, _ *http.Request) {
	a.renderLoginPage(w, nil)
	return
}

func (a *App) handlePostLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	acc := &model.AccountRequest{
		Username: r.PostFormValue("username"),
		Password: r.PostFormValue("password"),
	}
	tokenResp, response, err := a.APIClient.Login(ctx, acc)
	if err != nil {
		logger.Err(ctx, err)
		a.renderLoginPage(w, &model.ErrorPage{ErrorMessage: err.Error()})
		return
	}

	if util.StatusClientError(response.StatusCode) {
		logger.Info(ctx, "client error", "status", response.StatusCode, "message", response.Message)
		a.renderLoginPage(w, &model.ErrorPage{ErrorMessage: response.Message})
		return
	}

	if !util.StatusSuccess(response.StatusCode) {
		logger.Err(ctx, "error from server", "status code", response.StatusCode, "message", response.Message, "body", tokenResp)
		a.renderLoginPage(w, &model.ErrorPage{ErrorMessage: response.Message})
		return
	}

	a.setAccessTokenCookie(w, tokenResp)
	http.Redirect(w, r, "/profile/view", http.StatusFound)
}

func (a *App) renderLoginPage(w http.ResponseWriter, err *model.ErrorPage) {
	_ = a.Tmpl.Execute(w, templateLogin, model.LoginPage{
		ErrorPage:      err,
		LoginURI:       a.Config.LoginURI,
		GoogleClientID: a.Config.GoogleClientID,
	})
}

func (a *App) renderSignupPage(wr io.Writer, err *model.ErrorPage) {
	_ = a.Tmpl.Execute(wr, templateSignup, model.SignupPage{
		ErrorPage:      err,
		LoginURI:       a.Config.LoginURI,
		GoogleClientID: a.Config.GoogleClientID,
	})
}

func (a *App) handleGetSignup(w http.ResponseWriter, _ *http.Request) {
	a.renderSignupPage(w, nil)
}

func (a *App) handlePostSignup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	acc := &model.AccountRequest{
		Username: r.PostFormValue("username"),
		Password: r.PostFormValue("password"),
	}
	tokenResp, response, err := a.APIClient.Signup(ctx, acc)
	if err != nil {
		logger.Err(ctx, err)
		a.renderSignupPage(w, &model.ErrorPage{ErrorMessage: err.Error()})
		return
	}

	if !util.StatusSuccess(response.StatusCode) {
		logger.Err(ctx, "error from server", "response", fmt.Sprintf("%+v", response))
		a.renderSignupPage(w, &model.ErrorPage{ErrorMessage: response.Message})
		return
	}

	a.setAccessTokenCookie(w, tokenResp)
	http.Redirect(w, r, "/profile/edit", http.StatusFound)
}

func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie(cookieKeyAccessToken)
	if err != nil {
		logger.Info(ctx, "cannot find access token, redirect to login anyway...")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	accessToken := cookie.Value

	cookie.MaxAge = -1
	cookie.Value = ""
	http.SetCookie(w, cookie)

	if err = a.APIClient.Logout(ctx, accessToken); err != nil {
		logger.Err(ctx, err) // logout error should not interrupt logout behaviour
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *App) handleAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenResp, resp, err := a.APIClient.ForwardGoogleAuth(ctx, r)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	if !util.StatusSuccess(resp.StatusCode) {
		logger.Err(ctx, "error authentication from server", "response", *resp)
		_ = util.SendJSON(ctx, w, resp.StatusCode, "error authentication from server", nil)
		return
	}

	a.setAccessTokenCookie(w, tokenResp)
	http.Redirect(w, r, "/profile/view", http.StatusFound)
}

func (a *App) setAccessTokenCookie(w http.ResponseWriter, tokenResp *model.TokenResponse) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieKeyAccessToken,
		Value:    tokenResp.AccessToken,
		Secure:   true,
		HttpOnly: true,
	})
}
