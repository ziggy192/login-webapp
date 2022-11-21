package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/auth"
	"bitbucket.org/ziggy192/ng_lu/src/api/model"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"encoding/json"
	"net/http"
)

func (a *App) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := auth.GetUsername(ctx)
	profile, err := a.DBStores.Profile.FindProfileByUserName(ctx, username)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	if profile == nil {
		_ = util.SendJSON(ctx, w, http.StatusNotFound, "profile not found for username "+username, nil)
		return
	}

	_ = util.SendJSON(ctx, w, http.StatusOK, "get profile successfully", profile)
}

func (a *App) handleSaveProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := auth.GetUsername(r.Context())
	profile := &model.Profile{}
	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendJSON(ctx, w, http.StatusBadRequest, "unable to decode json profile", nil)
		return
	}
	profile.AccountID = username
	inserted, err := a.DBStores.Profile.InsertOrUpdateProfile(ctx, profile)
	if err != nil {
		logger.Err(ctx, err)
		_ = util.SendError(ctx, w, err)
		return
	}

	if inserted != 0 {
		profile.ID = inserted
	}

	_ = util.SendJSON(r.Context(), w, http.StatusOK, "save profile successfully", profile)
}
