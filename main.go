package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."
	// const filepathAssets = "./assets" // this along with 'http.StripePrefix' are no necesary!

	// create a new ServeMux (HTTP request multiplexer) to route incoming requests
	mux := http.NewServeMux()

	// build and run fileserver that serves `index.html` file from root at port 8080
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	// serve files from the assets directory
	// mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(filepathAssets))))

	// create and Initialize http.Server Struct (as Pointer)
	// memory efficiency: Only the pointer (which is small) is passed around,
	// not a copy of the entire struct instance
	server := &http.Server{
		// Addr:    ":" + port, // Binds to all interfaces
		Addr: "localhost:" + port, // Binds only to localhost
		Handler: mux, // this tells the server to route incoming requests using the ServeMux
	}

	log.Printf("...⚡ Serving files from %s on port: %s\n", filepathRoot, port)
	// log.Printf("...⚡ Serving on port: %s\n", port)
	log.Println("...💀 Stop the server using `Ctrl + C`")

	// start the server, use 'ListenAndServe' method,
	// using 'ListenAndServe', the main function blocks until the server is shut down
	// log.Fatal: logs the error returned by ListenAndServe and exits the program immediately
	log.Fatal(server.ListenAndServe())
}
