package domain

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password_hash"`
}

type ApiUser struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Passwords string `json:"password_hash"`
}
