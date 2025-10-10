package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/http2"

	"cloudrun/internal/httpapi"
)

func main() {
	addr := ":" + getEnv("PORT", "8080")
	srv := httpapi.NewServer()

	server := &http.Server{
		Addr:              addr,
		Handler:           srv.Router(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// HTTP/2 habilitado (sem TLS em Cloud Run é gerenciado pelo proxy)
	_ = http2.ConfigureServer(server, &http2.Server{})

	log.Printf("🚀 listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
