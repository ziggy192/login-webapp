package util

import (
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"bytes"
	"context"
	"io"
	"net/http"
)

const HeaderXRequestID = "X-REQUEST-ID"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body []byte
		if r.Body != nil {
			body, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		logger.Info(r.Context(), "request", r.RemoteAddr, r.Method, r.URL, string(body))
		next.ServeHTTP(w, r)
	})
}

// RequestIDMiddleware add/generate X-Request-ID value to the request context, and response's header for propagation
// and save request-id to context variable
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := r.Header.Get(HeaderXRequestID)
		var wrappedCtx context.Context
		if len(requestID) == 0 {
			wrappedCtx, requestID = logger.GenRequestID(ctx)
			r.Header.Set(HeaderXRequestID, requestID)
		} else {
			wrappedCtx = logger.SaveRequestID(ctx, requestID)
		}

		r = r.WithContext(wrappedCtx)
		w.Header().Set(HeaderXRequestID, requestID)

		next.ServeHTTP(w, r)
	})
}
