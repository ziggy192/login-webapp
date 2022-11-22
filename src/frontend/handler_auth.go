package frontend

import (
	"bitbucket.org/ziggy192/ng_lu/src/frontend/model"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	body, err := json.Marshal(acc)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}
	request, err := http.NewRequest(http.MethodPost, a.Config.APIRoot+pathLogin, bytes.NewBuffer(body))
	request.Header.Set(HeaderContentType, "application/json")
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}
	requestID := logger.GetRequestID(ctx)
	if len(requestID) > 0 {
		request.Header.Set(util.HeaderXRequestID, requestID)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	var tokenResp = &model.TokenResponse{}
	response, err := ParseResponse(ctx, resp, &tokenResp)
	if err != nil {
		logger.Err(ctx, err)
		a.renderLoginPage(w, &model.ErrorPage{ErrorMessage: err.Error()})
		return
	}

	logger.Info(ctx, "response", fmt.Sprintf("%+v", response))

	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		logger.Info(ctx, "unauthorized", "status", response.StatusCode, "message", response.Message)
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
	body, err := json.Marshal(acc)
	if err != nil {
		logger.Err(ctx, err)
		a.renderSignupPage(w, &model.ErrorPage{ErrorMessage: err.Error()})
		return
	}
	request, err := http.NewRequest(http.MethodPost, a.Config.APIRoot+pathSignup, bytes.NewBuffer(body))
	request.Header.Set(HeaderContentType, "application/json")
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}
	requestID := logger.GetRequestID(ctx)
	if len(requestID) > 0 {
		request.Header.Set(util.HeaderXRequestID, requestID)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	var tokenResp = &model.TokenResponse{}
	response, err := ParseResponse(ctx, resp, &tokenResp)
	if err != nil {
		logger.Err(ctx, err)
		a.renderSignupPage(w, &model.ErrorPage{ErrorMessage: err.Error()})
		return
	}

	logger.Info(ctx, "response", fmt.Sprintf("%+v", response))

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

	request, err := http.NewRequest(http.MethodPost, a.Config.APIRoot+pathLogout, nil)
	requestID := logger.GetRequestID(ctx)
	if len(requestID) > 0 {
		request.Header.Set(util.HeaderXRequestID, requestID)
	}

	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}
	request.Header.Set(headerAuthorization, "Bearer "+accessToken)
	logger.Info(ctx, "request", request)
	resp, err := http.DefaultClient.Do(request) // ignore response because logout error
	if err != nil {
		logger.Err(ctx, err) // error should not interrupt logout behavior
	}
	logger.Info(ctx, "response", resp)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *App) handleAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
	}

	request, err := http.NewRequest(http.MethodPost, a.Config.APIRoot+pathLoginGoogle, strings.NewReader(r.PostForm.Encode()))
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}
	requestID := logger.GetRequestID(ctx)
	if len(requestID) > 0 {
		request.Header.Set(util.HeaderXRequestID, requestID)
	}
	request.Header.Add(HeaderContentType, "application/x-www-form-urlencoded")

	csrfTokenCookie, err := r.Cookie("g_csrf_token")
	if err != nil {
		_ = util.SendJSON(ctx, w, 400, "no CSRF token in Cookie", nil)
		return
	}

	request.AddCookie(&http.Cookie{
		Name:   csrfTokenCookie.Name,
		Value:  csrfTokenCookie.Value,
		MaxAge: csrfTokenCookie.MaxAge,
	})
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}
	var tokenResp = &model.TokenResponse{}
	resp, err := ParseResponse(ctx, response, &tokenResp)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}
	logger.Info(ctx, "response", fmt.Sprintf("%+v", resp))

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
