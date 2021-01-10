package database

import (
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom-errors"
	"made.by.jst10/outfit7/hancock/cmd/structs"
)

func dbUserCreateTableIfNot() *custom_errors.CustomError {
	_, err := dbExec("CREATE TABLE  IF NOT EXISTS  users (" +
		"id int primary key auto_increment," +
		"created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"username varchar(128) NOT NULL," +
		"role int NOT NULL," +
		"password text NOT NULL," +
		"salt varchar(128) NOT NULL," +
		"CONSTRAINT users_username_key UNIQUE (username));")
	if err != nil {
		return err.AST("db create user table")
	} else {
		return nil
	}

}

func dbUserCreate(user *structs.User) (*structs.User, *custom_errors.CustomError) {
	res, err := dbExec("INSERT INTO users ("+
		"created_at, "+
		"updated_at, "+
		"username, "+
		"role, "+
		"password, "+
		"salt) values(CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?,?,?,?)",
		user.Username, user.Role, user.Password, user.Salt)
	if err != nil {
		return nil, err.AST("insert into user table")
	}
	id, err := dbLastInsertedId(res)
	if err != nil {
		return nil, err.AST("insert into user table")
	}
	return dbUserGetUserById(int(id))
}

func dbUserAll() ([]structs.User, *custom_errors.CustomError) {
	users := make([]structs.User, 0)
	results, err := dbQuery("SELECT * FROM users")
	if err != nil {
		return nil, err.AST("select from user table")
	}
	for results.Next() {
		var user structs.User
		err = dbScanRows(results,
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

func dbUserDeleteAll() *custom_errors.CustomError {
	_, err := dbExec("DELETE FROM users")
	if err != nil {
		return err.AST("delete from user table")
	} else {
		return nil
	}

}

func dbUserDeleteById(id int) *custom_errors.CustomError {
	_, err := dbExec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return err.AST("delete from user table by id")
	} else {
		return nil
	}

}

func dbUserGetUserById(id int) (*structs.User, *custom_errors.CustomError) {
	return _getOneUser("SELECT * FROM users WHERE id = ?", id)
}

func dbUserGetUserByUsername(username string) (*structs.User, *custom_errors.CustomError) {
	return _getOneUser("SELECT * FROM users WHERE username = ?", username)
}

func _getOneUser(query string, args ...interface{}) (*structs.User, *custom_errors.CustomError) {
	var user structs.User
	row := dbQueryRow(query, args...)
	err:=dbScanRow(row,
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Username,
		&user.Role,
		&user.Password,
		&user.Salt)
	if err != nil {
		return nil, err.AST("get one user")
	}
	return &user, nil
}
