package structs

type TokenData struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
	Role      int    `json:"role"`
	SessionId string `json:"session_id"`
}