package main

import (
	"net/http"
)

func main() {
	// create a new ServeMux (HTTP request multiplexer) to route incoming requests
	mux := http.NewServeMux()

	// create and Initialize http.Server Struct (as Pointer)
	// memory efficiency: Only the pointer (which is small) is passed around,
	// not a copy of the entire struct instance
	server := &http.Server{
		Addr: "localhost:8080",
		Handler: mux, // this tells the server to route incoming requests using the ServeMux
	}

	// start the server, use 'ListenAndServe' method
	server.ListenAndServe()
}
