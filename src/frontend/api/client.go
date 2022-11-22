package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/frontend/model"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bitbucket.org/ziggy192/ng_lu/src/util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	pathLogin       = "/login"
	pathLoginGoogle = "/login_google"
	pathSignup      = "/signup"
	pathLogout      = "/logout"
	pathProfile     = "/profile"

	HeaderAuthorization    = "Authorization"
	HeaderContentType      = "Content-Type"
	ContentTypeJSON        = "application/json"
	ContentTypeEncodedForm = "application/x-www-form-urlencoded"
)

// Client for backend REST API endpoints
type Client struct {
	apiRoot string
}

func NewClient(apiRoot string) *Client {
	return &Client{apiRoot: apiRoot}
}

func (c *Client) Login(ctx context.Context, acc *model.AccountRequest) (*model.TokenResponse, *Response, error) {
	body, err := json.Marshal(acc)
	if err != nil {
		logger.Err(ctx, err)
		return nil, nil, err
	}
	request, err := http.NewRequest(http.MethodPost, c.apiRoot+pathLogin, bytes.NewBuffer(body))
	request.Header.Set(HeaderContentType, "application/json")
	if err != nil {
		logger.Err(ctx, err)
		return nil, nil, err
	}

	requestID := logger.GetRequestID(ctx)
	if len(requestID) > 0 {
		request.Header.Set(util.HeaderXRequestID, requestID)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Err(ctx, err)
		return nil, nil, err
	}

	var tokenResp = &model.TokenResponse{}
	response, err := ParseResponse(ctx, resp, &tokenResp)
	if err != nil {
		logger.Err(ctx, err)
		return nil, response, err
	}

	logger.Info(ctx, "response", fmt.Sprintf("%+v", response))
	return tokenResp, response, nil
}

func (c *Client) Signup(ctx context.Context, acc *model.AccountRequest) (*model.TokenResponse, *Response, error) {
	body, err := json.Marshal(acc)
	if err != nil {
		logger.Err(ctx, err)
		return nil, nil, err
	}

	request, err := http.NewRequest(http.MethodPost, c.apiRoot+pathSignup, bytes.NewBuffer(body))
	request.Header.Set(HeaderContentType, ContentTypeJSON)
	if err != nil {
		logger.Err(ctx, err)
		return nil, nil, err
	}
	requestID := logger.GetRequestID(ctx)
	if len(requestID) > 0 {
		request.Header.Set(util.HeaderXRequestID, requestID)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Err(ctx, err)
		return nil, nil, err
	}

	var tokenResp = &model.TokenResponse{}
	response, err := ParseResponse(ctx, resp, &tokenResp)
	if err != nil {
		logger.Err(ctx, err)
		return nil, response, err
	}

	logger.Info(ctx, "response", fmt.Sprintf("%+v", response))
	return tokenResp, response, nil
}

func (c *Client) FetchProfile(ctx context.Context, accessToken string) (model.Profile, *Response, error) {
	request, err := http.NewRequest(http.MethodGet, c.apiRoot+pathProfile, nil)
	if err != nil {
		logger.Err(ctx, err)
		return model.Profile{}, nil, err
	}

	requestID := logger.GetRequestID(ctx)
	if len(requestID) > 0 {
		request.Header.Set(util.HeaderXRequestID, requestID)
	}
	request.Header.Set(HeaderAuthorization, "Bearer "+accessToken)

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
	logger.Info(ctx, "response", fmt.Sprintf("%+v", resp))
	return p, resp, nil
}

func (c *Client) PutProfile(ctx context.Context, accessToken string, profile model.Profile) (model.Profile, *Response, error) {
	body, err := json.Marshal(profile)
	if err != nil {
		logger.Err(ctx, err)
		return model.Profile{}, nil, err
	}

	request, err := http.NewRequest(http.MethodPut, c.apiRoot+pathProfile, bytes.NewBuffer(body))
	if err != nil {
		logger.Err(ctx, err)
		return model.Profile{}, nil, err
	}

	requestID := logger.GetRequestID(ctx)
	if len(requestID) > 0 {
		request.Header.Set(util.HeaderXRequestID, requestID)
	}
	request.Header.Set(HeaderAuthorization, "Bearer "+accessToken)

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
	logger.Info(ctx, "response", fmt.Sprintf("%+v", resp))
	return p, resp, nil
}

func (c *Client) Logout(ctx context.Context, accessToken string) error {
	request, err := http.NewRequest(http.MethodPost, c.apiRoot+pathLogout, nil)
	requestID := logger.GetRequestID(ctx)
	if len(requestID) > 0 {
		request.Header.Set(util.HeaderXRequestID, requestID)
	}

	if err != nil {
		logger.Err(ctx, err)
		return err
	}
	request.Header.Set(HeaderAuthorization, "Bearer "+accessToken)
	logger.Info(ctx, "request", request)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Err(ctx, err) // error should not interrupt logout behavior
	}
	logger.Info(ctx, "response", resp)
	return nil
}

func (c *Client) ForwardGoogleAuth(ctx context.Context, r *http.Request) (*model.TokenResponse, *Response, error) {
	err := r.ParseForm()
	if err != nil {
		logger.Err(ctx, err)
		return nil, nil, err
	}

	request, err := http.NewRequest(http.MethodPost, c.apiRoot+pathLoginGoogle, strings.NewReader(r.PostForm.Encode()))
	if err != nil {
		logger.Err(ctx, err)
		return nil, nil, err
	}
	requestID := logger.GetRequestID(ctx)
	if len(requestID) > 0 {
		request.Header.Set(util.HeaderXRequestID, requestID)
	}
	request.Header.Add(HeaderContentType, ContentTypeEncodedForm)

	csrfTokenCookie, err := r.Cookie("g_csrf_token")
	if err != nil {
		err = fmt.Errorf("no CSRF token in Cookie: %w", err)
		logger.Err(ctx, err)
		return nil, nil, err
	}

	request.AddCookie(&http.Cookie{
		Name:   csrfTokenCookie.Name,
		Value:  csrfTokenCookie.Value,
		MaxAge: csrfTokenCookie.MaxAge,
	})
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Err(ctx, err)
		return nil, nil, err
	}

	var tokenResp = &model.TokenResponse{}
	resp, err := ParseResponse(ctx, response, &tokenResp)
	if err != nil {
		logger.Err(ctx, err)
		return nil, nil, err
	}

	logger.Info(ctx, "response", fmt.Sprintf("%+v", resp))
	return tokenResp, resp, nil
}
