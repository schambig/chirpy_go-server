package main

import (
	"encoding/json"
	"net/http"
)

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
