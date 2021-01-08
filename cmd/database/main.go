package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

type User struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      int    `json:"role"`
	Password  string `json:"password"`
	Salt      string `json:"salt"`
}

type Session struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	UserId    string `json:"user_id"`
}

type Version struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
}

type Sdk struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type App struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Type struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func dbExec(query string) bool {
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
func createTables() {
	success := dbExec("CREATE TABLE IF NOT EXISTS versions(" +
		"id int primary key auto_increment," +
		"created_at TIMESTAMP default CURRENT_TIMESTAMP);")
	success = success && dbExec("CREATE TABLE IF NOT EXISTS countries("+
		"id int primary key,"+
		"name text);")
	success = success && dbExec("CREATE TABLE IF NOT EXISTS sdks("+
		"id int primary key,"+
		"name text);")
	success = success && dbExec("CREATE TABLE IF NOT EXISTS apps("+
		"id int primary key,"+
		"name text);")
	success = success && dbExec("CREATE TABLE IF NOT EXISTS performances("+
		"id int primary key auto_increment,"+
		"sdk int, " +
		"type int, "+
		"cmd int, "+
		"country int, "+
		"score int)")

	success = success && dbExec("CREATE TABLE  IF NOT EXISTS  users ("+
		"id int primary key auto_increment,"+
		"created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
		"updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
		"username varchar(100) NOT NULL,"+
		"first_name varchar(100) NOT NULL,"+
		"last_name varchar(100) NOT NULL,"+
		"email varchar(100) NOT NULL,"+
		"role int NOT NULL,"+
		"password text NOT NULL,"+
		"salt varchar(1024) NOT NULL,"+
		"CONSTRAINT users_email_key UNIQUE (email),"+
		"CONSTRAINT users_username_key UNIQUE (username));")

	success = success && dbExec("CREATE TABLE IF NOT EXISTS sessions ("+
		"id int primary key auto_increment,"+
		"created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
		"updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
		"user_id int NOT NULL);")

	if success {
		fmt.Println("All tables were successfully created.")
	}

}

func GetLatestVersion() (*Version, error) {
	fmt.Println("getla")
	var version Version
	results, err := db.Query("SELECT id, created_at FROM versions ORDER BY id DESC LIMIT 1;")
	if err != nil {
		return nil, err
	}
	fmt.Println("Gavesomething")
	for results.Next() {
		err = results.Scan(&version.ID, &version.CreatedAt)
		if err != nil {
			return nil, err
		}
	}
	return &version, nil
}

func createNewVersion() (*Version, error) {
	success := dbExec("INSERT INTO versions (created_at) VALUES (CURRENT_TIMESTAMP);")
	if success {
		return GetLatestVersion()
	} else {
		return nil, errors.New("Error at inseritng into versions")
	}
}

func InitDatabase() {
	var err error
	db, err = sql.Open("mysql",
		"root:root@tcp(127.0.0.1:3306)/hancock")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB connection has successfully initialized")
	createTables()
	GetLatestVersion()
	res, _ := createNewVersion()
	fmt.Println(res)
	//defer db.Close()
}
func a() {

}
