package main

import (
	"log"
	"net/http"
	"fmt"
)

// struct to hold any stateful (in-memory data)
type apiConfig struct {
	fileserverHits int
}

// middleware method to increment hits and call next handler
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits +=1
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hits: %d\n", cfg.fileserverHits)
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
}

// Handler function for /healthz
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	// create a new ServeMux (HTTP request multiplexer) to route incoming requests
	mux := http.NewServeMux()

	// intanciate from struct
	apiCfg := &apiConfig{}

	// use a variable to avoid long line
	fileServerHandler := http.FileServer(http.Dir(filepathRoot))
	// serve files from the root directory under the /app/* path, wrap fileserver with middleware
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServerHandler)))

	// endpoint registers
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerReset)

	server := &http.Server{
		// Addr:    ":" + port, // binds to all interfaces
		Addr: "localhost:" + port, // binds only to localhost
		Handler: mux, // this tells the server to route incoming requests using the ServeMux
	}

	log.Printf("...⚡ Serving files from %s on port: %s\n", filepathRoot, port)
	log.Println("...💀 Stop the server using `Ctrl + C`")

	log.Fatal(server.ListenAndServe())
}
