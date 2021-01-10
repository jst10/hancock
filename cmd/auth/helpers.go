package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom-errors"
	"made.by.jst10/outfit7/hancock/cmd/structs"
	"time"
)

var jwtKey = []byte("tokenKey")
var jwrtKey = []byte("refreshTokenKey")

func createSalt() (string, *custom_errors.CustomError) {
	var salt = make([]byte, 64)
	_, err := rand.Read(salt[:])
	if err != nil {
		return "", custom_errors.GetErrorCreateSalt(err)
	}
	return base64.URLEncoding.EncodeToString(salt), nil
}
func createHash(password string, salt string) string {
	var passwordBytes = []byte(password)
	var saltBytes = []byte(salt)
	var sha512Hasher = sha512.New()
	passwordBytes = append(passwordBytes, saltBytes...)
	sha512Hasher.Write(passwordBytes)
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	var base64EncodedPasswordHash = base64.URLEncoding.EncodeToString(hashedPasswordBytes)
	return base64EncodedPasswordHash
}

func doPasswordsMatch(user *structs.User, enteredPassword string) bool {
	var enteredPasswordHash = createHash(enteredPassword, user.Salt)
	return user.Password == enteredPasswordHash
}

func createJWT(tokenData *structs.TokenData) (string, *time.Time, *custom_errors.CustomError) {
	expirationTime := time.Now().Add(5 * time.Minute)
	tokenData.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", nil, custom_errors.GetErrorCreateJWT(err, "JWT")
	}
	return tokenString, &expirationTime, nil
}
func verifyJWT(token string) (*structs.TokenData, *custom_errors.CustomError) {
	tokenData := &structs.TokenData{}
	tkn, err := jwt.ParseWithClaims(token, tokenData, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, custom_errors.GetErrorVerifyJWT(err, "JWT")
	}
	if !tkn.Valid {
		return nil, custom_errors.GetInvalidJWT("JWT")
	}
	return tokenData, nil
}
func createJWRT(tokenData *structs.TokenData) (string, *time.Time, *custom_errors.CustomError) {
	expirationTime := time.Now().Add(365 * 24 * time.Hour)
	tokenData.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)
	tokenString, err := token.SignedString(jwrtKey)
	if err != nil {
		return "", nil, custom_errors.GetErrorCreateJWT(err, "JWRT")
	}
	return tokenString, &expirationTime, nil
}
func verifyJWRT(token string) (*structs.TokenData, *custom_errors.CustomError) {
	tokenData := &structs.TokenData{}
	tkn, err := jwt.ParseWithClaims(token, tokenData, func(token *jwt.Token) (interface{}, error) {
		return jwrtKey, nil
	})
	if err != nil {
		return nil, custom_errors.GetErrorVerifyJWT(err, "JWRT")
	}
	if !tkn.Valid {
		return nil, custom_errors.GetInvalidJWT("JWRT")
	}
	return tokenData, nil
}
