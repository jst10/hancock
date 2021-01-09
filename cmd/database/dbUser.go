package database

import (
	"made.by.jst10/outfit7/hancock/cmd/structs"
)

func dbUserCreateTableIfNot() error {
	_, err := db.Exec("CREATE TABLE  IF NOT EXISTS  users (" +
		"id int primary key auto_increment," +
		"created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"username varchar(128) NOT NULL," +
		"role int NOT NULL," +
		"password text NOT NULL," +
		"salt varchar(128) NOT NULL," +
		"CONSTRAINT users_username_key UNIQUE (username));")
	return err
}

func dbUserCreate(user structs.User) error {
	_, err := db.Exec("INSERT INTO users ("+
		"created_at, "+
		"updated_at, "+
		"username, "+
		"role, "+
		"password, "+
		"salt) values(CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?,?,?,?)",
		user.Username, user.Role, user.Password, user.Salt)
	return err
}

func dbUserAll() ([]structs.User, error) {
	users := make([]structs.User, 0)
	results, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var user structs.User
		err = results.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Username,
			&user.Role,
			&user.Password,
			&user.Salt)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func dbUserDeleteAll() error {
	_, err := db.Exec("DELETE FROM users")
	return err
}

func dbUserDeleteById(id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id=?", id)
	return err
}

func dbUserGetUserById(id int) (*structs.User, error) {
	return _getOneUser("SELECT * FROM users WHERE id = ?", id)
}

func dbUserGetUserByUsername(username string) (*structs.User, error) {
	return _getOneUser("SELECT * FROM users WHERE username = ?", username)
}

func _getOneUser(query string, args ...interface{}) (*structs.User, error) {
	var user structs.User
	err := db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Username,
		&user.Role,
		&user.Password,
		&user.Salt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
