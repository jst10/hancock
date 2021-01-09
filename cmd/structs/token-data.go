package structs
import "github.com/dgrijalva/jwt-go"
type TokenData struct {
	UserId        int    `json:"userId"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
	Role      int    `json:"role"`
	SessionId int    `json:"session_id"`
	jwt.StandardClaims
}
