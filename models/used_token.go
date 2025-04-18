package models

import "time"

type UsedToken struct {
	ID        int       `json:"id"`
	UserGuid  string    `json:"user_guid"`
	TokenHash string    `json:"token_hash"`
	UsedAt    time.Time `json:"used_at"`
}
