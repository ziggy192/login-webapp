package frontend

import (
	"bitbucket.org/ziggy192/ng_lu/src/frontend/model"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"encoding/json"
	"net/http"
	"strings"
)

const HeaderContentType = "Content-Type"

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == http.MethodGet {
		p := struct {
			LoginURI       string
			GoogleClientID string
		}{
			LoginURI:       a.Config.LoginURI,
			GoogleClientID: a.Config.GoogleClientID,
		}
		_ = a.Tmpl.ExecuteTemplate(w, "login.html", p)
		return
	}

	err := r.ParseForm()
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}
	var acc model.Account
	err = a.SchemaDecoder.Decode(&acc, r.PostForm)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	// todo do something with acc
	logger.Info(ctx, acc)

	// todo try login with account

	// returns fails if not success

	http.Redirect(w, r, "/profile/view", http.StatusFound)
}

func (a *App) handleSignup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == http.MethodGet {
		_ = a.Tmpl.ExecuteTemplate(w, "signup.html", nil)
		return
	}

	err := r.ParseForm()
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}
	var acc model.Account
	err = a.SchemaDecoder.Decode(&acc, r.PostForm)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	// todo do something with acc
	logger.Info(ctx, acc)

	// todo try signup with account

	// returns fails if not success

	http.Redirect(w, r, "/profile/edit", http.StatusFound)
}

func (a *App) handleProfileView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	// todo get profile from token

	cookie, err := r.Cookie(cookieKeyAccessToken)
	if err != nil {
		logger.Info(ctx, "cannot find \"access_token\" in cookie, redirect to login")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	accessToken := cookie.Value

	// todo separate package -> service
	request, err := http.NewRequest(http.MethodGet, a.Config.APIRoot+pathProfile, nil)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}
	requestID := logger.GetRequestID(ctx)
	if len(requestID) > 0 {
		request.Header.Set(util.HeaderXRequestID, requestID)
	}
	request.Header.Set(headerAuthorization, "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		logger.Info(ctx, "unauthorized, redirect to login", "status", resp.StatusCode)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if !util.StatusSuccess(resp.StatusCode) {
		logger.Err(ctx, "error from server", "status code", resp.StatusCode, "body", util.ReadBody(resp.Body))
		_ = util.SendJSON(ctx, w, resp.StatusCode, "error from server", nil)
		return
	}

	var baseResponse model.BaseResponse
	err = json.NewDecoder(resp.Body).Decode(&baseResponse)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	var p model.Profile
	err = baseResponse.UnmarshalData(&p)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
	}

	_ = a.Tmpl.ExecuteTemplate(w, "profile_view.html", p)
}

func (a *App) handleProfileEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
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
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}
	var p model.Profile
	err = a.SchemaDecoder.Decode(&p, r.PostForm)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	// todo do something with p
	logger.Info(ctx, p)

	// todo save p to database

	// todo returns fails if not success

	http.Redirect(w, r, "/profile/view", http.StatusFound)

	_ = a.Tmpl.ExecuteTemplate(w, "profile_edit.html", p)
}

func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// todo call api to logout
	cookie, err := r.Cookie(cookieKeyAccessToken)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	cookie.MaxAge = -1
	cookie.Value = ""
	http.SetCookie(w, cookie)

	accessToken := cookie.Value
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
	_, err = http.DefaultClient.Do(request) // ignore response because logout error
	if err != nil {
		logger.Err(ctx, err) // error should not interrupt logout behavior
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *App) handleAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
	}

	request, err := http.NewRequest(http.MethodPost, a.Config.APIRoot+loginGooglePath, strings.NewReader(r.PostForm.Encode()))
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
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	if !util.StatusSuccess(resp.StatusCode) {
		logger.Err(ctx, "error authentication from server", "status code", resp.StatusCode, "body", util.ReadBody(resp.Body))
		_ = util.SendJSON(ctx, w, resp.StatusCode, "error authentication from server", nil)
		return
	}

	var baseResponse model.BaseResponse
	if err = json.NewDecoder(resp.Body).Decode(&baseResponse); err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	var tokenResp = &model.TokenResponse{}
	if err = baseResponse.UnmarshalData(&tokenResp); err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  cookieKeyAccessToken,
		Value: tokenResp.AccessToken,
		//Secure:   true,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/profile/view", http.StatusFound)
}
