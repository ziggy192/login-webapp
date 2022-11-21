package logger

import (
	"context"
	"github.com/google/uuid"
	"strings"
)

type ctxKey string

var (
	ctxKeyRequestID ctxKey = "request_id"
)

// GenRequestID generates random request id and save to the context variable
func GenRequestID(ctx context.Context) (context.Context, string) {
	requestID := NewRequestID(ctx)
	return SaveRequestID(ctx, requestID), requestID
}

// NewRequestID generates random request id
func NewRequestID(ctx context.Context) string {
	genUUID, err := uuid.NewRandom()
	if err != nil {
		Err(ctx, err)
		return ""
	}
	return strings.ReplaceAll(genUUID.String(), "-", "")
}

// GetRequestID returns value of key either from context
func GetRequestID(ctx context.Context) string {
	if value, ok := ctx.Value(ctxKeyRequestID).(string); ok {
		return value
	}
	return ""
}

// SaveRequestID saves the request id to context variable and custom context store
// Return a copy of parent in which the value associated with key is
func SaveRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, requestID)
}
