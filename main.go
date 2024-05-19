package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	// create a new ServeMux (HTTP request multiplexer) to route incoming requests
	mux := http.NewServeMux()

	// serve files from the root directory under the /app/* path
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	// register the /healthz endpoint
	mux.HandleFunc("/healthz", handleReadiness)


	server := &http.Server{
		// Addr:    ":" + port, // Binds to all interfaces
		Addr: "localhost:" + port, // Binds only to localhost
		Handler: mux, // this tells the server to route incoming requests using the ServeMux
	}

	log.Printf("...⚡ Serving files from %s on port: %s\n", filepathRoot, port)
	log.Println("...💀 Stop the server using `Ctrl + C`")

	log.Fatal(server.ListenAndServe())
}

// Handler function for /healthz
func handleReadiness(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
