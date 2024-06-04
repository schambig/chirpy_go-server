package database

import (
	"encoding/json"
	"errors"
	"sync"
	"os"
)

type DB struct {
	path string
	mu *sync.RWMutex // pointer, to use a single/shared inst. of the mutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"` 
}

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`

}

// constructor function (struct instantiator), NOT a receiver method 
func NewDB(path string) (DB, error) {
	db := &DB{
		path: path,
		mu: &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

// receiver methods (use existing struct instance)
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.MarshalIndent(dbStructure, "", " ")
	if err != nil {
		return err
	}

	err := os.WriteFile(db.path, dat, 644)
	if err != nil {
		return err
	}
	return nil	
}
