package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"made.by.jst10/outfit7/hancock/cmd/structs"
)

func createSalt() (string, error) {
	var salt = make([]byte, 64)
	_, err := rand.Read(salt[:])
	if err != nil {
		return "", err
	}
	return string(salt), nil
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


