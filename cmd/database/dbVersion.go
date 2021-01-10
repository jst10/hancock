package database

import custom_errors "made.by.jst10/outfit7/hancock/cmd/custom-errors"

func dbVersionCreateTableIfNot()  *custom_errors.CustomError {
	_, err := dbExec("CREATE TABLE IF NOT EXISTS versions (" +
		"id int primary key auto_increment," +
		"created_at TIMESTAMP default CURRENT_TIMESTAMP," +
		"db_index int NOT NULL" +
		");")
	if err != nil {
		return err.AST("create version table")
	} else {
		return nil
	}
}

func dbVersionCreate(version *Version) (*Version, *custom_errors.CustomError) {
	res, err := dbExec("INSERT INTO versions ("+
		"created_at, db_index) "+
		"values(CURRENT_TIMESTAMP,?)",
		version.DbIndex)
	if err != nil {
		return nil, err.AST("create version")
	}
	id, err := dbLastInsertedId(res)
	if err != nil {
		return nil, err.AST("create version")
	}
	return dbVersionGetVersionById(int(id))
}

func dbVersionAll() ([]Version, *custom_errors.CustomError) {
	versions := make([]Version, 0)
	results, err := dbQuery("SELECT * FROM versions")
	if err != nil {
		return nil, err.AST("select from version table")
	}
	for results.Next() {
		var version Version
		err = dbScanRows(results,
			&version.ID,
			&version.CreatedAt,
			&version.DbIndex)
		if err != nil {
			continue
		}
		versions = append(versions, version)
	}
	return versions, nil
}

func dbVersionDeleteAll() *custom_errors.CustomError {
	_, err := dbExec("DELETE FROM versions")
	if err != nil {
		return err.AST("delete from version table")
	} else {
		return nil
	}
}

func dbVersionDeleteById(id int) *custom_errors.CustomError {
	_, err := dbExec("DELETE FROM versions WHERE id=?", id)
	if err != nil {
		return err.AST("delete by id from version table")
	} else {
		return nil
	}
}

func dbVersionGetVersionById(id int) (*Version, *custom_errors.CustomError) {
	return _getOneVersion("SELECT * FROM versions WHERE id = ?", id)
}

func dbVersionGetLatest() (*Version, *custom_errors.CustomError) {
	return _getOneVersion("SELECT * FROM versions ORDER BY id DESC LIMIT 1;")
}

func _getOneVersion(query string, args ...interface{}) (*Version, *custom_errors.CustomError) {
	var version Version
	row := dbQueryRow(query, args...)
	err:=dbScanRow(row,
		&version.ID,
		&version.CreatedAt,
		&version.DbIndex)
	if err != nil {
		return nil, err.AST("get one version")
	}
	return &version, nil
}
