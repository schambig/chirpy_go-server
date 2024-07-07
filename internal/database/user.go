package database

import (
	"errors"
)

type User struct {
	ID int `json:"id"`
	Password string `json:"password"`
	Email string `json:"email"`
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbStructure.Users) + 1
	user := User{
		ID: id,
		Password: password,
		Email: email,
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

	usrs, err := db.GetUsers()
	if err != nil {
		return User{}, err
	}

	for _, user := range usrs {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("Email not found")
}
