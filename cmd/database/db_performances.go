package database

import (
	"fmt"
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom_errors"
	"strconv"
)

const dbPerformanceTablePrefix = "performances"

func dbPerformanceCreateTablesIfNot(tableIndex int) *custom_errors.CustomError {
	tableName := dbPerformanceTablePrefix + strconv.Itoa(tableIndex)
	_, err := dbExec("CREATE TABLE IF NOT EXISTS " + tableName + "(" +
		"id int primary key," +
		"ad_type int, " +
		"country int, " +
		"app int, " +
		"sdk int, " +
		"score int)")
	if err != nil {
		return err.AST("create performance table")
	} else {
		return nil
	}
}

func dbPerformanceCreate(tableIndex int, performance *Performance) *custom_errors.CustomError {
	tableName := dbPerformanceTablePrefix + strconv.Itoa(tableIndex)
	_, err := dbExec("INSERT INTO "+tableName+" ("+
		"id, ad_type, country, app, sdk, score ) "+
		"values(?,?,?,?,?,?)",
		performance.ID, performance.AdType, performance.Country, performance.App, performance.Sdk, performance.Score)
	if err != nil {
		return err.AST("insert into performance table")
	} else {
		return nil
	}
}

func dbPerformanceCreateBulk(tableIndex int, performances []Performance) *custom_errors.CustomError {
	tableName := dbPerformanceTablePrefix + strconv.Itoa(tableIndex)
	//since are only integer there are no problems with sql injection
	limit := 30000
	for i := 0; i < len(performances); i += limit {
		query := "INSERT INTO " + tableName + " (id, ad_type, country, app, sdk, score ) VALUES"
		for j := i; j < min(i+limit, len(performances)); j++ {
			performance := performances[j]
			query = query + fmt.Sprintf("(%d, %d,%d, %d,%d, %d),", performance.ID, performance.AdType, performance.Country, performance.App, performance.Sdk, performance.Score)
		}
		query = query[:len(query)-1] + ";"
		_, err := dbExec(query)
		if err != nil {
			return err.AST("insert into performance table")
		}
	}
	return nil
}

func dbPerformanceAll(tableIndex int) ([]Performance, *custom_errors.CustomError) {
	tableName := dbPerformanceTablePrefix + strconv.Itoa(tableIndex)
	items := make([]Performance, 0)
	results, err := dbQuery("SELECT * FROM " + tableName + "  ORDER BY score DESC, sdk ASC")
	if err != nil {
		return nil, err.AST("select from performance table")
	}
	for results.Next() {
		var item Performance
		err = dbScanRows(results,
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

func dbPerformanceDeleteAll(tableIndex int) *custom_errors.CustomError {
	tableName := dbPerformanceTablePrefix + strconv.Itoa(tableIndex)
	_, err := dbExec("DELETE FROM " + tableName)
	if err != nil {
		return err.AST("delete from user table")
	} else {
		return nil
	}
}
