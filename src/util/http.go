package util

import (
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type BaseResponse struct {
	Message string `json:"message"`
	Time    string `json:"time"`
	Data    any    `json:"data"`
}

// SendJSON encodes data as JSON object and returns it to client
func SendJSON(ctx context.Context, w http.ResponseWriter, statusCode int, message string, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := BaseResponse{
		Message: message,
		Data:    data,
		Time:    time.Now().Format(time.RFC3339),
	}

	body, err := json.Marshal(resp)
	if err != nil {
		logger.Err(ctx, "cannot marshal response data:", err)
		return err
	}

	_, err = w.Write(body)
	if err != nil {
		logger.Err(ctx, "cannot write response body:", err)
		return err
	}

	return nil
}

// SendError sends internal error response to client
func SendError(ctx context.Context, w http.ResponseWriter, err error) error {
	return SendJSON(ctx, w, http.StatusInternalServerError, err.Error(), nil)
}

func StatusSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func StatusClientError(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}
