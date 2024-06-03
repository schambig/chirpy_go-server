package database

import (
	"encoding/json"
	"errors"
	"sync"
	"os"
)

type DB struct {
	path string
	mu sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"` 
}

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`

}

func (db *DB) createDB(path string)