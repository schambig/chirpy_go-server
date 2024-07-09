package database

import (
	"errors"
)

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
	HashedPassword string `json:"password"`
}

var ErrAlreadyExists = errors.New("Already exists")

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email, hashedPassword string) (User, error) {
	/* if _, err := db.GetUserByEmail(email); err != nil {
		if !errors.Is(err, ErrNotExist) {
			return User{}, ErrAlreadyExists
		}
	} */
	
	// same as above but shorter
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID: id,
		Email: email,
		HashedPassword: hashedPassword,
	}
	dbStructure.Users[id] = user
	
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetUsers returns all users in the database
func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	
	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}
