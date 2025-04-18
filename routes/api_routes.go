package routes

import (
	"med/db"
	"med/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func Init(r *mux.Router, manager *db.Manager) {

}

func Token(r *mux.Router, manager *db.Manager) {
	r.HandleFunc("/token/{guid}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAccessRefreshToken(manager, w, r)
	}).Methods(http.MethodPost)
}

func RefreshToken(r *mux.Router, manager *db.Manager) {
	r.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		handlers.RefreshAccessToken(manager, w, r)
	}).Methods(http.MethodPost)
}
