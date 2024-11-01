package types

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminPayload struct {
	UserID string `json:"user_id"`
}
