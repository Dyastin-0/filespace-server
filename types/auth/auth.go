package auth

type Body struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	AccessToken string   `json:"accessToken"`
	Email       string   `json:"user"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
}
