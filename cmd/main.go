package main

import (
	"fmt"
	"made.by.jst10/outfit7/hancock/cmd/database"
)

func performCheck() {
	fmt.Println("Perform: Check")
}

func main() {
	//performCheck()
	//handleRequests()
	database.InitDatabase()
}
