package middlewares

import (
	"fmt"
	"net/http"
)

const cacheMaxAge = 1800 // the DB basically never changes

func CacheHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d, stale-if-error=120", cacheMaxAge))
		next.ServeHTTP(w, r)
	})
}
