package database

import "strconv"

const dbPerformanceTablePrefix = "performances"

func dbPerformanceCreateTablesIfNot(tableIndex int) error {
	tableName := dbPerformanceTablePrefix + strconv.Itoa(tableIndex)
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(" +
		"id int primary key," +
		"sdk int, " +
		"type int, " +
		"app int, " +
		"country int, " +
		"score int)")
	return err
}

func dbPerformanceCreate(tableIndex int, performance *Performance) error {
	tableName := dbPerformanceTablePrefix + strconv.Itoa(tableIndex)
	_, err := db.Exec("INSERT INTO "+tableName+" ("+
		"id, ad_type, country, app, sdk,  score ) "+
		"values(?,?,?,?,?,?)",
		performance.ID, performance.AdType, performance.Country, performance.App, performance.Sdk, performance.Score)
	return err
}

func dbPerformanceAll(tableIndex int) ([]Performance, error) {
	tableName := dbPerformanceTablePrefix + strconv.Itoa(tableIndex)
	items := make([]Performance, 0)
	results, err := db.Query("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var item Performance
		err = results.Scan(
			&item.ID,
			&item.AdType,
			&item.Country,
			&item.App,
			&item.Sdk,
			&item.Score)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func dbPerformanceDeleteAll(tableIndex int) error {
	tableName := dbPerformanceTablePrefix + strconv.Itoa(tableIndex)
	_, err := db.Exec("DELETE FROM " + tableName)
	return err
}
