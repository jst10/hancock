package auth

import (
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom-errors"
	"made.by.jst10/outfit7/hancock/cmd/database"
	"made.by.jst10/outfit7/hancock/cmd/structs"
)

func CreateUser(user *structs.User) (*structs.User, *custom_errors.CustomError) {
	var err *custom_errors.CustomError
	user.Salt, err = createSalt()
	if err != nil {
		return nil, err.AST("create salt")
	}
	user.Password = createHash(user.Password, user.Salt)
	return database.CreateUser(user)
}

func AuthenticateUser(authData *structs.AuthData) (*structs.TokenWrapper, *structs.TokenWrapper, *custom_errors.CustomError) {
	user, err := database.GetUserByUsername(authData.Username)
	if err != nil {
		return nil, nil, err.AST("authenticate user")
	}
	if !doPasswordsMatch(user, authData.Password) {
		return nil, nil, custom_errors.GetNotValidDataError("invalid auth data")
	}
	session, err := database.CreateSession(&structs.Session{UserId: user.ID})
	if err != nil {
		return nil, nil, err.AST("authenticate user")
	}
	tokenData := structs.TokenData{
		UserId:    user.ID,
		CreatedAt: user.CreatedAt,
		Username:  user.Username,
		Role:      user.Role,
		SessionId: session.ID}

	token, tokenExpiration, err := createJWT(&tokenData)
	if err != nil {
		return nil, nil, err.AST("authenticate user")
	}

	refreshToken, refreshTokenExpiration, err := createJWRT(&tokenData)
	if err != nil {
		return nil, nil, err.AST("authenticate user")
	}

	tokenWrapper := structs.TokenWrapper{User: user, Token: token, Expiration: tokenExpiration}
	refreshTokenWrapper := structs.TokenWrapper{User: user, Token: refreshToken, Expiration: refreshTokenExpiration}
	return &tokenWrapper, &refreshTokenWrapper, nil
}
func ReAuthenticateUser(sessionId int) (*structs.TokenWrapper, *custom_errors.CustomError) {
	session, err := database.GetSessionById(sessionId)
	if err != nil {
		return nil, err.AST("re authenticate user")
	}
	user, err := database.GetUserById(session.UserId)
	if err != nil {
		return nil, err.AST("re authenticate user")
	}
	tokenData := structs.TokenData{
		UserId:    user.ID,
		CreatedAt: user.CreatedAt,
		Username:  user.Username,
		Role:      user.Role,
		SessionId: session.ID}

	token, tokenExpiration, err := createJWT(&tokenData)
	if err != nil {
		return nil, err.AST("re authenticate user")
	}

	return &structs.TokenWrapper{User: user, Token: token, Expiration: tokenExpiration}, nil
}
func DeAuthenticateUser(userId int) *custom_errors.CustomError {
	user, err := database.GetUserById(userId)
	if err != nil {
		return err.AST("de authenticate user")
	}
	return database.DeleteUserSessions(user.ID)
}

func VerifyJWT(token string) (*structs.TokenData, *custom_errors.CustomError) {
	return verifyJWT(token)
}
func VerifyJWRT(token string) (*structs.TokenData, *custom_errors.CustomError) {
	return verifyJWRT(token)
}
