package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"sort"
	"log"
	"os"
)

// struct for the json body to expect
type validChirp struct {
	Body string `json:"body"`
}

// struct to return marshaled JSON
type returnChirp struct {
	Id int `json:"id"`
	Body string `json:"body"`
}

// struct to hold next id state (in-memory data)
type chirpId struct {
	nextID int
	mu sync.RWMutex
}

var chirpCounter = &chirpId{}

type DBStructure struct {
	Chirps map[int]returnChirp `json:"chirps"`
}

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

	newChirp := returnChirp{
		Id: id,
		Body: cleanedBody,
	}

	dbStructure := loadChirpsFromFile()
	dbStructure.Chirps[id] = newChirp

	err = writeToDatabaseFile(dbStructure)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save chirp")
		return
	}

	// respond with successful message if all went as expected
	respondWithJSON(w, http.StatusCreated, returnChirp{
		Id: id,
		Body: cleanedBody,
	})
}

func handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbStructure := loadChirpsFromFile()

	var chirps []returnChirp
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	respondWithJSON(w, http.StatusOK, chirps)
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

func writeToDatabaseFile(data interface{}) error {
	dat, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	
	return os.WriteFile("./database.json", dat, 0644)
}

func loadChirpsFromFile() DBStructure {
	file, err := os.ReadFile("./database.json")
	if err != nil {
		if os.IsNotExist(err) {
			return DBStructure{Chirps: make(map[int]returnChirp)}
		}
		log.Printf("Error reading database file: %v", err)
		return DBStructure{Chirps: make(map[int]returnChirp)}
	}

	var dbStruct DBStructure
	err = json.Unmarshal(file, &dbStruct)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return DBStructure{Chirps: make(map[int]returnChirp)}
	}
	return dbStruct
}
