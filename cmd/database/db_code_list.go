package database

import (
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom_errors"
	"strconv"
)
const dbAppTablePrefix = "apps"
const dbSdkTablePrefix = "sdks"
const dbCountryTablePrefix = "countries"

func dbSdkCreateTablesIfNot(tableIndex int) *custom_errors.CustomError {
	tableName := dbSdkTablePrefix + strconv.Itoa(tableIndex)
	return dbCodeListCreateTablesIfNot(tableName)
}

func dbSdkCreate(tableIndex int, item *CodeList) *custom_errors.CustomError {
	tableName := dbSdkTablePrefix + strconv.Itoa(tableIndex)
	return dbCodeListCreate(tableName, item)
}

func dbSdkAll(tableIndex int) ([]*CodeList, *custom_errors.CustomError) {
	tableName := dbSdkTablePrefix + strconv.Itoa(tableIndex)
	return dbCodeListAll(tableName)
}

func dbSdkDeleteAll(tableIndex int) *custom_errors.CustomError {
	tableName := dbSdkTablePrefix + strconv.Itoa(tableIndex)
	return dbCodeListDeleteAll(tableName)
}

func dbCountryCreateTablesIfNot(tableIndex int) *custom_errors.CustomError {
	tableName := dbCountryTablePrefix + strconv.Itoa(tableIndex)
	return dbCodeListCreateTablesIfNot(tableName)
}
func dbCountryCreate(tableIndex int, item *CodeList) *custom_errors.CustomError {
	tableName := dbCountryTablePrefix + strconv.Itoa(tableIndex)
	return dbCodeListCreate(tableName, item)
}
func dbCountryAll(tableIndex int) ([]*CodeList, *custom_errors.CustomError) {
	tableName := dbCountryTablePrefix + strconv.Itoa(tableIndex)
	return dbCodeListAll(tableName)
}
func dbCountryDeleteAll(tableIndex int) *custom_errors.CustomError {
	tableName := dbCountryTablePrefix + strconv.Itoa(tableIndex)
	return dbCodeListDeleteAll(tableName)
}

func dbAppCreateTablesIfNot(tableIndex int) *custom_errors.CustomError {
	tableName := dbAppTablePrefix + strconv.Itoa(tableIndex)
	return dbCodeListCreateTablesIfNot(tableName)
}
func dbAppCreate(tableIndex int, item *CodeList) *custom_errors.CustomError {
	tableName := dbAppTablePrefix + strconv.Itoa(tableIndex)
	return dbCodeListCreate(tableName, item)
}
func dbAppAll(tableIndex int) ([]*CodeList, *custom_errors.CustomError) {
	tableName := dbAppTablePrefix + strconv.Itoa(tableIndex)
	return dbCodeListAll(tableName)
}
func dbAppDeleteAll(tableIndex int) *custom_errors.CustomError {
	tableName := dbAppTablePrefix + strconv.Itoa(tableIndex)
	return dbCodeListDeleteAll(tableName)
}



func dbCodeListCreateTablesIfNot(tableName string) *custom_errors.CustomError {
	_, err := dbExec("CREATE TABLE IF NOT EXISTS " + tableName + "(" +
		"id int primary key," +
		"name text);")
	if err != nil {
		return err.AST("create codeList table")
	} else {
		return nil
	}
}

func dbCodeListCreate(tableName string, codeList *CodeList) *custom_errors.CustomError {
	_, err := dbExec("INSERT INTO "+tableName+" (id, name ) values(?,?)", codeList.ID, codeList.Name)
	if err != nil {
		return err.AST("insert into codeList table")
	} else {
		return nil
	}
}

func dbCodeListAll(tableName string) ([]*CodeList, *custom_errors.CustomError) {
	items := make([]*CodeList, 0)
	results, err := dbQuery("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err.AST("select from codeList table")
	}
	for results.Next() {
		var item CodeList
		err = dbScanRows(results,
			&item.ID,
			&item.Name)
		if err != nil {
			continue
		}
		items = append(items, &item)
	}
	return items, nil
}

func dbCodeListDeleteAll(tableName string) *custom_errors.CustomError {
	_, err := dbExec("DELETE FROM " + tableName)
	if err != nil {
		return err.AST("delete from codeList table")
	} else {
		return nil
	}
}
