package handlers

import (
	"med/db"
	"net/http"
)

func GetAccessRefreshToken(manager *db.Manager, w http.ResponseWriter, r *http.Request) {}

func RefreshAccessToken(manager *db.Manager, w http.ResponseWriter, r *http.Request) {}