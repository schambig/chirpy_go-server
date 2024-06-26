package main

import (
	"net/http"
	"sync"
	"log"

	"github.com/schambig/chirpy_go-server/internal/database"
)

// struct to hold any stateful (in-memory data)
type apiConfig struct {
	fileserverHits int
	mu sync.RWMutex // only need one RWMutex to handle both reads and writes

	DB *database.DB
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	// intanciate from struct
	apiCfg := apiConfig{
		DB: db, // just initialize DB field from struct
	}

	// create a new ServeMux (HTTP request multiplexer) to route incoming requests
	mux := http.NewServeMux()

	// use a variable to avoid long line
	fileServerHandler := http.FileServer(http.Dir(filepathRoot))
	// serve files from the root directory under the /app/* path, wrap fileserver with middleware
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServerHandler)))

	// endpoint registers
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpID)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUsers)

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
