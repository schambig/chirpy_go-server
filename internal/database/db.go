package database

import (
	"encoding/json"
	"errors"
	"sync"
	"os"
)

var ErrNotExist = errors.New("Resource does not exist")

type DB struct {
	path string
	mu *sync.RWMutex // pointer, to use a single/shared inst. of the mutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
}

// constructor function (struct instantiator), NOT a receiver method

// NewDB creates a new database connection
// and creates the database file if it doesn't exist 
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu: &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

// receiver methods (use existing struct instance)

// ensureDB creates a new database file if it doesn't exist

func (db *DB) ResetDB() error {
	err := os.Remove(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return db.ensureDB()
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	// must use os package for the app to check the OS for the database.json file
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users: map[int]User{},
	}
	return db.writeDB(dbStructure)
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.MarshalIndent(dbStructure, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0644)
	if err != nil {
		return err
	}
	
	return nil	
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()	

	// initialize dbStructure with a non-nil Chirps map to avoid runtime error
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users: map[int]User{},
	}
	dat, err := os.ReadFile(db.path)
	// since it won't create the database.json if error happens, `os` is optional
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}
