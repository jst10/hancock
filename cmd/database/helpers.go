package database

import (
	"database/sql"
	"fmt"
)

func dbExec(db *sql.DB,query string) bool {
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
