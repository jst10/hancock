package auth

import (
	"crypto/sha512"
	"encoding/base64"
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

//func doPasswordsMatch(hashedPassword, currPassword string, salt string) bool {
//	var saltBytes = []byte(salt)
//	var currPasswordHash = hashPassword(currPassword, saltBytes)
//	return hashedPassword == currPasswordHash
//}
