package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Response is a response data
type Response struct {
	StatusCode  int       `json:"status_code"`
	ContentType string    `json:"content_type"`
	Message     string    `json:"message"`
	Time        time.Time `json:"time"`
}

// parsedBody stores parsed json body
// Data should be a pointer
type parsedBody struct {
	Message string
	Time    time.Time
	Data    any
}

// EmptyResponse is an empty response
var EmptyResponse = &Response{}

// ParseResponse parses http response to struct
// outData should be a pointer
func ParseResponse(ctx context.Context, r *http.Response, outData interface{}) (*Response, error) {
	if r == nil {
		return EmptyResponse, errors.New("http response is nil")
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			logger.Err(ctx, err)
		}
	}()

	response := &Response{
		StatusCode:  r.StatusCode,
		ContentType: r.Header.Get(HeaderContentType),
	}

	body := parsedBody{Data: outData}

	contentType := r.Header.Get(HeaderContentType)
	if !strings.HasPrefix(contentType, ContentTypeJSON) {
		err := fmt.Errorf("body is not a json, status code = %d, content type = %s", r.StatusCode, contentType)
		logger.Err(ctx, err)
		return response, err
	}

	if err := parseResponseFromReader(ctx, r.Body, &body); err != nil {
		return response, err
	}

	response.Message = body.Message
	response.Time = body.Time

	return response, nil
}

func parseResponseFromReader(ctx context.Context, reader io.Reader, outData interface{}) error {
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()

	if err := decoder.Decode(outData); err != nil {
		logger.Err(ctx, err)
		return err
	}

	return nil
}
