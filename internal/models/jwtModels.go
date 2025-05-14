package models

type RefreshToken struct {
	Token string `json:"refresh_token"`
}

type AccessToken struct {
	Token string `json:"access_token"`
}
