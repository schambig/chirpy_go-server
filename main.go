package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"log"
)

// struct to hold any stateful (in-memory data)
type apiConfig struct {
	fileserverHits int
	mu sync.RWMutex // only need one RWMutex to handle both reads and writes
}

// struct for the json body to expect
type validChirp struct {
	Body string `json:"body"`
}

func handlerValidChirp(w http.ResponseWriter, r *http.Request) {
	var chirp validChirp

	// decode the json request body into the chirp variable
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)

	if err != nil {
		// respond with error if json decoding fails
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error":"Something went wrong when decoding JSON"})
		return
	}

	// check length of the chirp (Body field)
	if len(chirp.Body) > 140 {
		// respond with error if Body field exceeds length
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error":"Chirp is too long"})
		return		
	}
	
	// respond with successful message if all went as expected
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"valid":true})
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
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidChirp)

	server := &http.Server{
		// Addr:    ":" + port, // binds to all interfaces
		Addr: "localhost:" + port, // binds only to localhost
		Handler: mux, // this tells the server to route incoming requests using the ServeMux
	}

	log.Printf("...⚡ Serving files from %s on port: %s\n", filepathRoot, port)
	log.Println("...💀 Stop the server using `Ctrl + C`")

	log.Fatal(server.ListenAndServe())
}
