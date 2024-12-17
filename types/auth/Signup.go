package auth

type SignupBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}
