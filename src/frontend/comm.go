package frontend

import (
	"bitbucket.org/ziggy192/ng_lu/src/frontend/model"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

func (a *App) fetchProfile(ctx context.Context, accessToken string) (model.Profile, *Response, error) {
	request, err := http.NewRequest(http.MethodGet, a.Config.APIRoot+pathProfile, nil)
	if err != nil {
		logger.Err(ctx, err)
		return model.Profile{}, nil, err
	}

	requestID := logger.GetRequestID(ctx)
	if len(requestID) > 0 {
		request.Header.Set(util.HeaderXRequestID, requestID)
	}
	request.Header.Set(headerAuthorization, "Bearer "+accessToken)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Err(ctx, err)
		return model.Profile{}, nil, err
	}

	var p model.Profile
	resp, err := ParseResponse(ctx, response, &p)
	if err != nil {
		logger.Err(ctx, err)
		return model.Profile{}, nil, err
	}
	return p, resp, nil
}

func (a *App) putProfile(ctx context.Context, accessToken string, profile model.Profile) (model.Profile, *Response, error) {
	body, err := json.Marshal(profile)
	if err != nil {
		logger.Err(ctx, err)
		return model.Profile{}, nil, err
	}

	request, err := http.NewRequest(http.MethodPut, a.Config.APIRoot+pathProfile, bytes.NewBuffer(body))
	if err != nil {
		logger.Err(ctx, err)
		return model.Profile{}, nil, err
	}

	requestID := logger.GetRequestID(ctx)
	if len(requestID) > 0 {
		request.Header.Set(util.HeaderXRequestID, requestID)
	}
	request.Header.Set(headerAuthorization, "Bearer "+accessToken)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Err(ctx, err)
		return model.Profile{}, nil, err
	}

	var p model.Profile
	resp, err := ParseResponse(ctx, response, &p)
	if err != nil {
		logger.Err(ctx, err)
		return model.Profile{}, nil, err
	}
	return p, resp, nil
}
