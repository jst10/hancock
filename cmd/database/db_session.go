package database

import (
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom_errors"
	"made.by.jst10/outfit7/hancock/cmd/structs"
)

func dbSessionCreateTableIfNot()  *custom_errors.CustomError {
	_, err := dbExec("CREATE TABLE IF NOT EXISTS sessions (" +
		"id int primary key auto_increment," +
		"created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"user_id int NOT NULL);")
	if err != nil {
		return err.AST("create session table")
	} else {
		return nil
	}
}

func dbSessionCreate(session *structs.Session) (*structs.Session, *custom_errors.CustomError) {
	res, err := dbExec("INSERT INTO sessions ("+
		"created_at, "+
		"updated_at, "+
		"user_id "+
		") values(CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?)",
		session.UserId)
	if err != nil {
		return nil, err.AST("insert into session table")
	}
	id, err := dbLastInsertedId(res)
	if err != nil {
		return nil, err.AST("insert into session table")
	}
	return dbSessionGetSessionById(int(id))
}

func dbSessionAll() ([]structs.Session, *custom_errors.CustomError) {
	return _getListOfSession("SELECT * FROM sessions")
}

func dbSessionsByUserId(userId int) ([]structs.Session, *custom_errors.CustomError) {
	return _getListOfSession("SELECT * FROM sessions WHERE user_id=?", userId)
}

func dbSessionDeleteAll() *custom_errors.CustomError {
	_, err := dbExec("DELETE FROM sessions")
	if err != nil {
		return err.AST("delete from session table")
	} else {
		return nil
	}
}

func dbSessionDeleteByUserId(userId int) *custom_errors.CustomError {
	_, err := dbExec("DELETE FROM sessions WHERE user_id=?", userId)
	if err != nil {
		return err.AST("delete from session table by user id")
	} else {
		return nil
	}
}

func dbSessionDeleteById(id int) *custom_errors.CustomError {
	_, err := dbExec("DELETE FROM sessions WHERE id=?", id)
	if err != nil {
		return err.AST("delete from session table by id")
	} else {
		return nil
	}
}

func dbSessionGetSessionById(id int) (*structs.Session, *custom_errors.CustomError) {
	return _getOneSession("SELECT * FROM sessions WHERE id = ?", id)
}

func _getListOfSession(query string, args ...interface{}) ([]structs.Session, *custom_errors.CustomError) {
	sessions := make([]structs.Session, 0)
	results, err := dbQuery(query, args...)
	if err != nil {
		return nil, err.AST("get list of sessions")
	}
	for results.Next() {
		var session structs.Session
		err = dbScanRows(results,
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

func _getOneSession(query string, args ...interface{}) (*structs.Session, *custom_errors.CustomError) {
	var session structs.Session
	row := dbQueryRow(query, args...)
	err := dbScanRow(row,
		&session.ID,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.UserId)
	if err != nil {
		return nil, err.AST("get one session")
	}
	return &session, nil
}
