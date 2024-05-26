package main

import (
	"net/http"
)

// Handler function for /healthz
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte(http.StatusText(http.StatusOK)))
	w.Write([]byte("OK\n")) // simpler and better option than line above
}
