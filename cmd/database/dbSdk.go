package database

import "strconv"

const dbSdkTablePrefix = "sdks"

func dbSdkCreateTablesIfNot(tableIndex int) error {
	tableName := dbSdkTablePrefix + strconv.Itoa(tableIndex)
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(" +
		"id int primary key," +
		"name text);")
	return err
}
func dbSdkCreate(tableIndex int, sdk *Sdk) error {
	tableName := dbSdkTablePrefix + strconv.Itoa(tableIndex)
	_, err := db.Exec("INSERT INTO "+tableName+" (id, name ) values(?,?)", sdk.ID, sdk.Name)
	return err
}
func dbSdkAll(tableIndex int) ([]Sdk, error) {
	tableName := dbSdkTablePrefix + strconv.Itoa(tableIndex)
	items := make([]Sdk, 0)
	results, err := db.Query("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var item Sdk
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
func dbSdkDeleteAll(tableIndex int) error {
	tableName := dbSdkTablePrefix + strconv.Itoa(tableIndex)
	_, err := db.Exec("DELETE FROM " + tableName)
	return err
}
