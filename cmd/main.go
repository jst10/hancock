package main

import (
	"log"
	"made.by.jst10/outfit7/hancock/cmd/api"
	"made.by.jst10/outfit7/hancock/cmd/auth"
	"made.by.jst10/outfit7/hancock/cmd/config"
	"made.by.jst10/outfit7/hancock/cmd/constants"
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom_errors"
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

func insertDefaultUserInDBIfNot() *custom_errors.CustomError {
	_, err := database.GetUserByUsername(userUser.Username)
	if err != nil {
		_, err := auth.CreateUser(&userUser)
		if err != nil {
			return err.AST("insert default user in db if not")
		}
	}
	_, err = database.GetUserByUsername(userAdmin.Username)
	if err != nil {
		_, err := auth.CreateUser(&userAdmin)
		if err != nil {
			return err.AST("insert default user in db if not")
		}
	}
	return nil
}

func main() {
	appConfigs := &config.AppConfigs{}
	err := config.LoadConfig(appConfigs)
	if err != nil {
		log.Fatal(err)
	}
	err = database.InitDatabase(appConfigs.Db)
	if err != nil {
		log.Fatal(err)
	}
	err=insertDefaultUserInDBIfNot()
	if err != nil {
		log.Fatal(err)
	}
	api.StartApi(appConfigs.Api)
}
