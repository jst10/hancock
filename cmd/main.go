package main

import (
	"log"
	"made.by.jst10/outfit7/hancock/cmd/auth"
	"made.by.jst10/outfit7/hancock/cmd/constants"
	"made.by.jst10/outfit7/hancock/cmd/database"
	"made.by.jst10/outfit7/hancock/cmd/structs"
)

var userUser = structs.User{
	Username: "user",
	Password: "user",
	Role:     constants.UserRoleGuest,
}

var userAdmin = structs.User{
	Username: "admin",
	Password: "admin",
	Role:     constants.UserRoleAdmin,
}

func insertDefaultUserInDBIfNot() error {
	_, err := database.GetUserByUsername(userUser.Username)
	if err != nil {
		_, err := auth.CreateUser(&userUser)
		if err != nil {
			return err
		}
	}
	_, err = database.GetUserByUsername(userAdmin.Username)
	if err != nil {
		_, err := auth.CreateUser(&userAdmin)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	err=insertDefaultUserInDBIfNot()
	if err != nil {
		log.Fatal(err)
	}
	startApi()
}
