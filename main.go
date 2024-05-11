package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	// create a new ServeMux (HTTP request multiplexer) to route incoming requests
	mux := http.NewServeMux()

	// create and Initialize http.Server Struct (as Pointer)
	// memory efficiency: Only the pointer (which is small) is passed around,
	// not a copy of the entire struct instance
	server := &http.Server{
		Addr: "localhost:" + port, // // Binds only to localhost
		Handler: mux, // this tells the server to route incoming requests using the ServeMux
	}

	log.Printf("...⚡ Serving on port: %v\n", port)
	log.Println("...💀 Close the server with `Ctrl + C`")

	// start the server, use 'ListenAndServe' method,
	// log.Fatal: logs the error returned by ListenAndServe and exits the program immediately
	log.Fatal(server.ListenAndServe())
}

/* 
func main() {
	const port = "8080"

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port, // // Binds to all interfaces
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
 */
