package database

import (
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom-errors"
	"strconv"
)

const dbCountryTablePrefix = "countries"

func dbCountryCreateTablesIfNot(tableIndex int) *custom_errors.CustomError {
	tableName := dbCountryTablePrefix + strconv.Itoa(tableIndex)
	_, err := dbExec("CREATE TABLE IF NOT EXISTS " + tableName + "(" +
		"id int primary key," +
		"name text);")
	if err != nil {
		return err.AST("create country table")
	} else {
		return nil
	}
}
func dbCountryCreate(tableIndex int, country *Country) *custom_errors.CustomError {
	tableName := dbCountryTablePrefix + strconv.Itoa(tableIndex)
	_, err := dbExec("INSERT INTO "+tableName+" (id, name ) values(?,?)", country.ID, country.Name)
	if err != nil {
		return err.AST("insert into country table")
	} else {
		return nil
	}
}
func dbCountryAll(tableIndex int) ([]Country, *custom_errors.CustomError) {
	tableName := dbCountryTablePrefix + strconv.Itoa(tableIndex)
	items := make([]Country, 0)
	results, err := dbQuery("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err.AST("select from country table")
	}
	for results.Next() {
		var item Country
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
func dbCountryDeleteAll(tableIndex int) *custom_errors.CustomError {
	tableName := dbCountryTablePrefix + strconv.Itoa(tableIndex)
	_, err := dbExec("DELETE FROM " + tableName)
	if err != nil {
		return err.AST("delete from country table")
	} else {
		return nil
	}
}
