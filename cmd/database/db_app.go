package database

import (
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom_errors"
	"strconv"
)

const dbAppTablePrefix = "apps"

func dbAppCreateTablesIfNot(tableIndex int) *custom_errors.CustomError {
	tableName := dbAppTablePrefix + strconv.Itoa(tableIndex)
	_, err := dbExec("CREATE TABLE IF NOT EXISTS " + tableName + "(" +
		"id int primary key," +
		"name text);")
	if err != nil {
		return err.AST("create app table")
	} else {
		return nil
	}
}
func dbAppCreate(tableIndex int, app *App) *custom_errors.CustomError {
	tableName := dbAppTablePrefix + strconv.Itoa(tableIndex)
	_, err := dbExec("INSERT INTO "+tableName+" (id, name ) values(?,?)", app.ID, app.Name)
	if err != nil {
		return err.AST("insert into app table")
	} else {
		return nil
	}
}
func dbAppAll(tableIndex int) ([]App, *custom_errors.CustomError) {
	tableName := dbAppTablePrefix + strconv.Itoa(tableIndex)
	items := make([]App, 0)
	results, err := dbQuery("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err.AST("select from app table")
	}
	for results.Next() {
		var item App
		err = dbScanRows(results,
			&item.ID,
			&item.Name)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}
func dbAppDeleteAll(tableIndex int) *custom_errors.CustomError {
	tableName := dbAppTablePrefix + strconv.Itoa(tableIndex)
	_, err := dbExec("DELETE FROM " + tableName)
	if err != nil {
		return err.AST("delete from app table")
	} else {
		return nil
	}
}
