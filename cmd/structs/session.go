package structs

type Session struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	UserId    int `json:"user_id"`
}