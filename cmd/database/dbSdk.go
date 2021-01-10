package database

import (
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom-errors"
	"strconv"
)

const dbSdkTablePrefix = "sdks"

func dbSdkCreateTablesIfNot(tableIndex int) *custom_errors.CustomError {
	tableName := dbSdkTablePrefix + strconv.Itoa(tableIndex)
	_, err := dbExec("CREATE TABLE IF NOT EXISTS " + tableName + "(" +
		"id int primary key," +
		"name text);")
	if err != nil {
		return err.AST("create sdk table")
	} else {
		return nil
	}
}
func dbSdkCreate(tableIndex int, sdk *Sdk) *custom_errors.CustomError {
	tableName := dbSdkTablePrefix + strconv.Itoa(tableIndex)
	_, err := dbExec("INSERT INTO "+tableName+" (id, name ) values(?,?)", sdk.ID, sdk.Name)
	if err != nil {
		return err.AST("insert into sdk table")
	} else {
		return nil
	}
}
func dbSdkAll(tableIndex int) ([]Sdk, *custom_errors.CustomError) {
	tableName := dbSdkTablePrefix + strconv.Itoa(tableIndex)
	items := make([]Sdk, 0)
	results, err := dbQuery("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err.AST("select from sdk table")
	}
	for results.Next() {
		var item Sdk
		err =dbScanRows(results,
			&item.ID,
			&item.Name)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}
func dbSdkDeleteAll(tableIndex int) *custom_errors.CustomError {
	tableName := dbSdkTablePrefix + strconv.Itoa(tableIndex)
	_, err := dbExec("DELETE FROM " + tableName)
	if err != nil {
		return err.AST("delete from sdk table")
	} else {
		return nil
	}
}
