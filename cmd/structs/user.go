package structs

type User struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Username  string `json:"username"`
	Role      int    `json:"role"`
	Password  string `json:"password"`
	Salt      string `json:"salt"`
}