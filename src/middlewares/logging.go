package middlewares

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logEntry := logrus.WithField("remote", r.RemoteAddr).WithField("meth", r.Method)
		if r.Header.Get("User-Agent") != "" {
			logEntry = logEntry.WithField("User-Agent", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("X-Forwarded-For") != "" {
			logEntry = logEntry.WithField("X-Forwarded-For", r.Header.Get("X-Forwarded-For"))
		}
		if r.Header.Get("X-Real-IP") != "" {
			logEntry = logEntry.WithField("X-Real-IP", r.Header.Get("X-Real-IP"))
		}
		logEntry.Debugf("Request: %s", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
