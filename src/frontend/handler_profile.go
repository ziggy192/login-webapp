package frontend

import (
	"bitbucket.org/ziggy192/ng_lu/src/frontend/model"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"net/http"
)

func (a *App) handleProfileView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	cookie, err := r.Cookie(cookieKeyAccessToken)
	if err != nil {
		logger.Info(ctx, "cannot find \"access_token\" in cookie, redirect to login")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	accessToken := cookie.Value
	p, resp, err := a.APIClient.FetchProfile(ctx, accessToken)
	if err != nil {
		a.renderProfileViewpage(w, p, &model.ErrorPage{ErrorMessage: resp.Message})
		return
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		logger.Info(ctx, "unauthorized, redirect to login", "status", resp.StatusCode)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if !util.StatusSuccess(resp.StatusCode) && resp.StatusCode != http.StatusNotFound {
		logger.Err(ctx, "error from server", "response", *resp)
		a.renderProfileViewpage(w, p, &model.ErrorPage{ErrorMessage: resp.Message})
		return
	}

	a.renderProfileViewpage(w, p, nil)
}

func (a *App) renderProfileViewpage(w http.ResponseWriter, profile model.Profile, errorPage *model.ErrorPage) {
	p := model.ProfileViewPage{
		Profile:   profile,
		ErrorPage: errorPage,
	}
	_ = a.Tmpl.Execute(w, templateProfileView, p)
}

func (a *App) renderProfileEditPage(w http.ResponseWriter, profile model.Profile, errorPage *model.ErrorPage) {
	p := model.ProfileEditPage{
		Profile:   profile,
		ErrorPage: errorPage,
	}
	_ = a.Tmpl.Execute(w, templateProfileEdit, p)
}
func (a *App) handleGetProfileEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	cookie, err := r.Cookie(cookieKeyAccessToken)
	if err != nil {
		logger.Info(ctx, "cannot find \"access_token\" in cookie, redirect to login")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	accessToken := cookie.Value
	p, resp, err := a.APIClient.FetchProfile(ctx, accessToken)
	if err != nil {
		a.renderProfileEditPage(w, p, &model.ErrorPage{ErrorMessage: resp.Message})
		return
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		logger.Info(ctx, "unauthorized, redirect to login", "status", resp.StatusCode)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if resp.StatusCode == http.StatusNotFound {
		logger.Info(ctx, "profile not yet created", "response", *resp)
		a.renderProfileEditPage(w, p, nil)
		return
	}

	if !util.StatusSuccess(resp.StatusCode) {
		logger.Err(ctx, "error from server", "response", *resp)
		a.renderProfileEditPage(w, p, &model.ErrorPage{ErrorMessage: resp.Message})
		return
	}

	a.renderProfileEditPage(w, p, nil)
}

func (a *App) handlePostProfileEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	err := r.ParseForm()
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	var p = model.Profile{
		FullName: r.PostFormValue("full_name"),
		Phone:    r.PostFormValue("phone"),
		Email:    r.PostFormValue("email"),
	}

	cookie, err := r.Cookie(cookieKeyAccessToken)
	if err != nil {
		logger.Info(ctx, "cannot find \"access_token\" in cookie, redirect to login")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	accessToken := cookie.Value
	p, response, err := a.APIClient.PutProfile(ctx, accessToken, p)

	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		logger.Info(ctx, "unauthorized, redirect to login", "status", response.StatusCode)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if !util.StatusSuccess(response.StatusCode) {
		logger.Err(ctx, "error from server", "response", *response)
		a.renderProfileEditPage(w, p, &model.ErrorPage{ErrorMessage: response.Message})
		return
	}

	http.Redirect(w, r, "/profile/view", http.StatusFound)
}
