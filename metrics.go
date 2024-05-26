package main

import (
	"fmt"
	"net/http"
)

// middleware method to increment hits and call next handler
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.mu.Lock() // lock for writing
		defer cfg.mu.Unlock() // ensure unlock after writing

		cfg.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.mu.RLock() // lock for reading
	defer cfg.mu.RUnlock() // ensure unlock after reading	

	htmlTemplate := `
	<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	</html>	
	`
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, htmlTemplate, cfg.fileserverHits)
}
