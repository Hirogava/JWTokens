package models

type User struct {
	Guid string `json:"guid"`
	Ip string `json:"ip"`
	Email string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	AccessId string `json:"accessToken"`
	CreatedAt string `json:"createdAt"`
}