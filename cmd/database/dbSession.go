package database

import (
	"made.by.jst10/outfit7/hancock/cmd/structs"
)

func dbSessionCreateTableIfNot() error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS sessions (" +
		"id int primary key auto_increment," +
		"created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"user_id int NOT NULL);")
	return err
}

func dbSessionCreate(session structs.Session) error {
	_, err := db.Exec("INSERT INTO sessions ("+
		"created_at, "+
		"updated_at, "+
		"user_id "+
		") values(CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?)",
		session.UserId)
	return err
}

func dbSessionAll() ([]structs.Session, error) {
	return _getListOfSession("SELECT * FROM sessions")
}

func dbSessionsByUserId(userId int) ([]structs.Session, error) {
	return _getListOfSession("SELECT * FROM sessions WHERE user_id=?", userId)
}

func dbSessionDeleteAll() error {
	_, err := db.Exec("DELETE FROM sessions")
	return err
}

func dbSessionDeleteByUserId(userId int) error {
	_, err := db.Exec("DELETE FROM sessions WHERE user_id=?", userId)
	return err
}

func dbSessionDeleteById(id int) error {
	_, err := db.Exec("DELETE FROM sessions WHERE id=?", id)
	return err
}

func dbSessionGetSessionById(id int) (*structs.Session, error) {
	return _getOneSession("SELECT * FROM sessions WHERE id = ?", id)
}

func _getListOfSession(query string, args ...interface{}) ([]structs.Session, error) {
	sessions := make([]structs.Session, 0)
	results, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var session structs.Session
		err = results.Scan(
			&session.ID,
			&session.CreatedAt,
			&session.UpdatedAt,
			&session.UserId)
		if err != nil {
			continue
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func _getOneSession(query string, args ...interface{}) (*structs.Session, error) {
	var session structs.Session
	err := db.QueryRow(query, args...).Scan(
		&session.ID,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.UserId)
	if err != nil {
		return nil, err
	}
	return &session, nil
}
