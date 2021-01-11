package database

import (
	"database/sql"
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom_errors"
)

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func sqlOpen(driverName, dataSourceName string) (*sql.DB, *custom_errors.CustomError) {
	result, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, custom_errors.GetDbOpenError(err)
	} else {
		return result, nil
	}
}

func dbPing() *custom_errors.CustomError {
	err := db.Ping()
	if err != nil {
		return custom_errors.GetDbPingError(err)
	} else {
		return nil
	}
}

func dbExec(query string, args ...interface{}) (sql.Result, *custom_errors.CustomError) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return nil, custom_errors.GetDbExecError(err)
	} else {
		return result, nil
	}
}

func dbQuery(query string, args ...interface{}) (*sql.Rows, *custom_errors.CustomError) {
	result, err := db.Query(query, args...)
	if err != nil {
		return nil, custom_errors.GetDbQueryError(err)
	} else {
		return result, nil
	}
}

func dbQueryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func dbScanRow(rs *sql.Row, dest ...interface{}) *custom_errors.CustomError {
	err := rs.Scan(dest...)
	if err != nil {
		return custom_errors.GetDbScanError(err)
	} else {
		return nil
	}
}

func dbScanRows(rs *sql.Rows, dest ...interface{}) *custom_errors.CustomError {
	err := rs.Scan(dest...)
	if err != nil {
		return custom_errors.GetDbScanError(err)
	} else {
		return nil
	}
}

func dbLastInsertedId(res sql.Result) (int64, *custom_errors.CustomError) {
	id, err := res.LastInsertId()
	if err != nil {
		return -1, custom_errors.GetDbGetLastInsertedIdError(err)
	} else {
		return id, nil
	}
}
