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

func (db *DB) UpdateUser(userID int, email, hashedPassword string ) (User, error) {
	// for modifying the database, use Lock and Unlock
	// But since there are mutexes in WriteBD func we don't need mutexes here (to avoid a deadlock)
	//     The first lock in UpdateUser won't release until the func returns, but WriteDB tries to
	//     acquire the same lock, causing the app to wait indefinitely (a mutex deadlock)
	// db.mu.Lock()
	// defer db.mu.Unlock()
	
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// ensure the user exists
	// Go's map structure can return 2 values when accessing an element by key
	// - The value associated with the key (if it exists)
	// - A boolean indicating if the key was found
	user, exists := dbStructure.Users[userID]
	if !exists {
		return User{}, ErrNotExist
	}

	user.Email = email
	user.HashedPassword = hashedPassword
	dbStructure.Users[userID] = user

	// write the updated structure back to the database
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil	
}
