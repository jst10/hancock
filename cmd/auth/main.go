package auth

import (
	"errors"
	"fmt"
	"made.by.jst10/outfit7/hancock/cmd/database"
	"made.by.jst10/outfit7/hancock/cmd/structs"
)

func CreateUser(user *structs.User) error {
	var err error
	user.Salt, err = createSalt()
	if err != nil {
		return err
	}
	user.Password = createHash(user.Password, user.Salt)
	return database.CreateUser(user)
}

func AuthenticateUser(authData *structs.AuthData) error {
	user, err := database.GetUserByUsername(authData.Username)
	if err != nil {
		return err
	}
	if !doPasswordsMatch(user, authData.Password) {
		return errors.New("invalid credentials")
	}

	return nil
}
func ReAuthenticateUser(sessionId int) error {
	session, err := database.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	user, err := database.GetUserById(session.UserId)
	if err != nil {
		return err
	}
	fmt.Println(user)

	return nil
}
func DeAuthenticateUser(userId int) error {
	user, err := database.GetUserById(userId)
	if err != nil {
		return err
	}
	return database.DeleteUserSessions(user.ID)
}
