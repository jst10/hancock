package database

import "strconv"

const dbCountryTablePrefix = "countries"

func dbCountryCreateTablesIfNot(tableIndex int) error {
	tableName := dbCountryTablePrefix + strconv.Itoa(tableIndex)
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(" +
		"id int primary key," +
		"name text);")
	return err
}
func dbCountryCreate(tableIndex int, country *Country) error {
	tableName := dbCountryTablePrefix + strconv.Itoa(tableIndex)
	_, err := db.Exec("INSERT INTO "+tableName+" (id, name ) values(?,?)", country.ID, country.Name)
	return err
}
func dbCountryAll(tableIndex int) ([]Country, error) {
	tableName := dbCountryTablePrefix + strconv.Itoa(tableIndex)
	items := make([]Country, 0)
	results, err := db.Query("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var item Country
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
func dbCountryDeleteAll(tableIndex int) error {
	tableName := dbCountryTablePrefix + strconv.Itoa(tableIndex)
	_, err := db.Exec("DELETE FROM " + tableName)
	return err
}
