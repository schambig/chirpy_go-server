package main

import (
	"net/http"
	"sync"
	"log"
)

// struct to hold any stateful (in-memory data)
type apiConfig struct {
	fileserverHits int
	mu sync.RWMutex // only need one RWMutex to handle both reads and writes
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
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServerHandler)))

	// endpoint registers
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", handlerValidChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	server := &http.Server{
		// Addr:    ":" + port, // binds to all interfaces
		Addr: "localhost:" + port, // binds only to localhost
		Handler: mux, // this tells the server to route incoming requests using the ServeMux
	}

	log.Printf("...⚡ Serving files from %s on port: %s\n", filepathRoot, port)
	log.Println("...💀 Stop the server using `Ctrl + C`")

	log.Fatal(server.ListenAndServe())
}
