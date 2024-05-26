package main

import (
	"fmt"
	"net/http"
)

// middleware method to increment hits and call next handler
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits +=1
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hits: %d\n", cfg.fileserverHits)
	htmlTemplate := `
    <html>
    <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
    </body>
    </html>	
	`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, htmlTemplate, cfg.fileserverHits)
}
