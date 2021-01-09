package database

import "strconv"

const dbAppTablePrefix = "apps"

func dbAppCreateTablesIfNot(tableIndex int) error {
	tableName := dbAppTablePrefix + strconv.Itoa(tableIndex)
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(" +
		"id int primary key," +
		"name text);")
	return err
}
func dbAppCreate(tableIndex int, app *App) error {
	tableName := dbAppTablePrefix + strconv.Itoa(tableIndex)
	_, err := db.Exec("INSERT INTO "+tableName+" (id, name ) values(?,?)", app.ID, app.Name)
	return err
}
func dbAppAll(tableIndex int) ([]App, error) {
	tableName := dbAppTablePrefix + strconv.Itoa(tableIndex)
	items := make([]App, 0)
	results, err := db.Query("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var item App
		err = results.Scan(
			&item.ID,
			&item.Name)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}
func dbAppDeleteAll(tableIndex int) error {
	tableName := dbAppTablePrefix + strconv.Itoa(tableIndex)
	_, err := db.Exec("DELETE FROM " + tableName)
	return err
}
