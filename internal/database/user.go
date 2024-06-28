package database

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(dbStructure.Users) + 1
	user := User{
		ID: id,
		Email: email,
	}
	dbStructure.Users[id] = user
	
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

/* // GetUsers returns all users in the database
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

func (db *DB) GetUserByID(id int) (User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}
	return user, nil
} */
