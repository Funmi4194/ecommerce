package types

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Verification struct {
	Email string `json:"email"`
	Link  string `json:"link"`
}

type Confirmation struct {
	Token string `json:"token"`
}

type PasswordRecovery struct {
	Email    string `json:"email"`
	Link     string `json:"link"`
	Token    string `json:"token"`
	Password string `json:"password"`
}
