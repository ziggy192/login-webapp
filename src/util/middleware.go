package util

import (
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		runTime := time.Since(startTime)
		logger.Info("request", r.RemoteAddr, r.Method, r.URL, runTime, ReadBody(r.Body))
	})
}
