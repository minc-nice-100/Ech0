package caphttp

import (
	"net"
	"net/http"
	"strings"
)

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]any{
					"success": false,
					"error":   "Internal server error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request, ipHeader string) string {
	if ipHeader != "" {
		v := r.Header.Get(ipHeader)
		if v != "" {
			parts := strings.Split(v, ",")
			return strings.TrimSpace(parts[0])
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
