package main

import (
	"encoding/json"
	"net/http"
	"log"
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
		respondWithError(w, http.StatusInternalServerError, "Something went wrong when decoding JSON")
		return
	}
	
	// check length of the chirp (Body field)
	const maxChirpLength = 140
	if len(chirp.Body) > maxChirpLength {
		// respond with error if Body field exceeds length
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return		
	}
	
	// respond with successful message if all went as expected
	respondWithJSON(w, http.StatusOK, map[string]bool{"valid":true})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code >= 500 {
		log.Printf("Server error 5XX: %s,", msg)
	} else if code >= 400 {
		log.Printf("Client error 4XX: %s,", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}
	
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error while marshalling JSON: %s,", err)
		w.WriteHeader(500) // marshalling error is a server error
		return
	}
	log.Printf("Successfully marshalled!")
	w.WriteHeader(code)
	w.Write(dat)
}
