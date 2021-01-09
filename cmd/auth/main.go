package auth

import (
	"errors"
	"made.by.jst10/outfit7/hancock/cmd/database"
	"made.by.jst10/outfit7/hancock/cmd/structs"
)

func CreateUser(user *structs.User) (*structs.User, error) {
	var err error
	user.Salt, err = createSalt()
	if err != nil {
		return nil, err
	}
	user.Password = createHash(user.Password, user.Salt)
	return database.CreateUser(user)
}

func AuthenticateUser(authData *structs.AuthData) (*structs.TokenWrapper, *structs.TokenWrapper, error) {
	user, err := database.GetUserByUsername(authData.Username)
	if err != nil {
		return nil, nil, err
	}
	if !doPasswordsMatch(user, authData.Password) {
		return nil, nil, errors.New("invalid credentials")
	}
	session, err := database.CreateSession(&structs.Session{UserId: user.ID})
	if err != nil {
		return nil, nil, err
	}
	tokenData := structs.TokenData{
		UserId:    user.ID,
		CreatedAt: user.CreatedAt,
		Username:  user.Username,
		Role:      user.Role,
		SessionId: session.ID}

	token, tokenExpiration, err := createJWT(&tokenData)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, refreshTokenExpiration, err := createJWRT(&tokenData)
	if err != nil {
		return nil, nil, err
	}

	tokenWrapper := structs.TokenWrapper{User: user, Token: token, Expiration: tokenExpiration}
	refreshTokenWrapper := structs.TokenWrapper{User: user, Token: refreshToken, Expiration: refreshTokenExpiration}
	return &tokenWrapper, &refreshTokenWrapper, nil
}
func ReAuthenticateUser(sessionId int) (*structs.TokenWrapper, error) {
	session, err := database.GetSessionById(sessionId)
	if err != nil {
		return nil, err
	}
	user, err := database.GetUserById(session.UserId)
	if err != nil {
		return nil, err
	}
	tokenData := structs.TokenData{
		UserId:    user.ID,
		CreatedAt: user.CreatedAt,
		Username:  user.Username,
		Role:      user.Role,
		SessionId: session.ID}

	token, tokenExpiration, err := createJWT(&tokenData)
	if err != nil {
		return nil, err
	}

	return &structs.TokenWrapper{User: user, Token: token, Expiration: tokenExpiration}, nil
}
func DeAuthenticateUser(userId int) error {
	user, err := database.GetUserById(userId)
	if err != nil {
		return err
	}
	return database.DeleteUserSessions(user.ID)
}

func VerifyJWT(token string) (*structs.TokenData, error) {
	return verifyJWT(token)
}
func VerifyJWRT(token string) (*structs.TokenData, error) {
	return verifyJWRT(token)
}
