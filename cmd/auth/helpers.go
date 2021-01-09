package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"made.by.jst10/outfit7/hancock/cmd/structs"
	"time"
)

var jwtKey = "tokenKey"
var jwrtKey = "refreshTokenKey"

func createSalt() (string, error) {
	var salt = make([]byte, 64)
	_, err := rand.Read(salt[:])
	if err != nil {
		return "", err
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

func createJWT(tokenData *structs.TokenData) (string, *time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	tokenData.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", nil, err
	}
	return tokenString, &expirationTime, nil
}
func verifyJWT(token string) (*structs.TokenData, error) {
	tokenData := &structs.TokenData{}
	tkn, err := jwt.ParseWithClaims(token, tokenData, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}
	return tokenData, nil
}
func createJWRT(tokenData *structs.TokenData) (string, *time.Time, error) {
	expirationTime := time.Now().Add(365 * 24 * time.Hour)
	tokenData.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)
	tokenString, err := token.SignedString(jwrtKey)
	if err != nil {
		return "", nil, err
	}
	return tokenString, &expirationTime, nil
}
func verifyJWRT(token string) (*structs.TokenData, error) {
	tokenData := &structs.TokenData{}
	tkn, err := jwt.ParseWithClaims(token, tokenData, func(token *jwt.Token) (interface{}, error) {
		return jwrtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}
	return tokenData, nil
}
