package database

func dbVersionCreateTableIfNot() error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS versions (" +
		"id int primary key auto_increment," +
		"created_at TIMESTAMP default CURRENT_TIMESTAMP," +
		"db_index int NOT NULL" +
		");")
	return err
}

func dbVersionCreate(version *Version) (*Version, error) {
	res, err := db.Exec("INSERT INTO versions ("+
		"created_at, db_index) "+
		"values(CURRENT_TIMESTAMP,?)",
		version.DbIndex)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return dbVersionGetVersionById(int(id))
}

func dbVersionAll() ([]Version, error) {
	versions := make([]Version, 0)
	results, err := db.Query("SELECT * FROM versions")
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var version Version
		err = results.Scan(
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

func dbVersionDeleteAll() error {
	_, err := db.Exec("DELETE FROM versions")
	return err
}

func dbVersionDeleteById(id int) error {
	_, err := db.Exec("DELETE FROM versions WHERE id=?", id)
	return err
}

func dbVersionGetVersionById(id int) (*Version, error) {
	return _getOneVersion("SELECT * FROM versions WHERE id = ?", id)
}

func dbVersionGetLatest() (*Version, error) {
	return _getOneVersion("SELECT * FROM versions ORDER BY id DESC LIMIT 1;")
}

func _getOneVersion(query string, args ...interface{}) (*Version, error) {
	var version Version
	err := db.QueryRow(query, args...).Scan(
		&version.ID,
		&version.CreatedAt,
		&version.DbIndex)
	if err != nil {
		return nil, err
	}
	return &version, nil
}
