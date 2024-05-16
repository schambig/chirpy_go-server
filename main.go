package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	server := &http.Server{
		// Addr:    ":" + port, // Binds to all interfaces
		Addr: "localhost:" + port, // Binds only to localhost
		Handler: mux, // this tells the server to route incoming requests using the ServeMux
	}

	log.Printf("...⚡ Serving files from %s on port: %s\n", filepathRoot, port)
	log.Println("...💀 Stop the server using `Ctrl + C`")
	log.Fatal(server.ListenAndServe())
}
