package types

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
	User        User
}
