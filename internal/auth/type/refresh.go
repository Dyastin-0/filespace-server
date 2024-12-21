package types

type User struct {
	Username string   `bson:"username"`
	Email    string   `bson:"email"`
	Roles    []string `bson:"roles"`
	ImageURL string   `bson:"profileImageURL"`
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
	User        User
}
