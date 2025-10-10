package limiter

import (
	"net"
	"net/http"
	"strings"
)

func getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func RateLimitMiddleware(l *Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("API_KEY")
			ip := getClientIP(r)

			var allowed bool
			var err error

			if token != "" {
				allowed, err = l.AllowToken(token)
			} else {
				allowed, err = l.AllowIP(ip)
			}

			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
