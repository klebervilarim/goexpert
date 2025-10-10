package main

import (
	"fmt"
	"net/http"
	"rate-limiter/limiter"
)

func main() {
	l := limiter.NewLimiter()

	mux := http.NewServeMux()
	mux.Handle("/test", limiter.RateLimitMiddleware(l)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("request ok"))
	})))

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}
