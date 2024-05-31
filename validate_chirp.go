package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

// struct for the json body to expect
type validChirp struct {
	Body string `json:"body"`
}

// struct to return marshaled JSON
type returnChirp struct {
	Body string `json:"body"`
	Id int `json:"id"`
}

// struct to hold next id state (in-memory data)
type chirpId struct {
	nextID int
	mu sync.RWMutex
}

var chirpCounter = &chirpId{}

func handlerValidChirp(w http.ResponseWriter, r *http.Request) {
	var chirp validChirp

	// decode the json request body into the chirp variable
	defer r.Body.Close()
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

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	cleanedBody := replaceProfaneWords(chirp.Body, badWords)
	id := chirpCounter.getID()

	// respond with successful message if all went as expected
	respondWithJSON(w, http.StatusCreated, returnChirp{
		Body: cleanedBody,
		Id: id,
	})
}

func replaceProfaneWords(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func (ci *chirpId) getID() int {
	chirpCounter.mu.Lock()
	defer chirpCounter.mu.Unlock()

	ci.nextID += 1
	return ci.nextID
}
