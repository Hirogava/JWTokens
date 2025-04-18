package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"med/db"
	"med/db/users"
	"med/models"
	"med/services/tokens"
	"net/http"
)

func GetAccessRefreshToken(manager *db.Manager, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Println("Неправильный Content-Type")
		http.Error(w, "Content-Type должен быть application/json", http.StatusBadRequest)
		return
	}

	var requestData struct {
		ID    string `json:"id"`
		IP    string `json:"ip"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Println("Ошибка при декодировании запроса:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User

	user, err := users.GetUserByGuid(requestData.ID, manager)
	if err == sql.ErrNoRows {
		user, err = users.CreateUser(manager, requestData.IP, requestData.Email)
		if err != nil {
			log.Println("Ошибка при создании пользователя:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		log.Println("Ошибка при получении пользователя:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := updateUser(manager, user)
	if err != nil {
		log.Println("Ошибка при обновлении пользователя:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func RefreshAccessToken(manager *db.Manager, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Println("Неправильный Content-Type")
		http.Error(w, "Content-Type должен быть application/json", http.StatusBadRequest)
		return
	}

	var requestData struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		ID           string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Println("Ошибка при декодировании запроса:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := users.GetUserByGuid(requestData.ID, manager)
	if err != nil {
		log.Println("Пользователь не найден", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, refreshTokenHash, err := tokens.UpdateTokens(manager, user, requestData.RefreshToken, requestData.AccessToken)
	if err != nil {
		log.Println("Ошибка при обновлении токенов:", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if _, err := users.UpdateUserRefreshToken(manager, user.Guid, refreshTokenHash); err != nil {
		log.Println("Ошибка при обновлении токена пользователя:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		log.Println("Ошибка при сериализации ответа:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func updateUser(manager *db.Manager, user models.User) ([]byte, error) {
	accesTokenString, refreshToken, refreshTokenHash, err := tokens.GenerateAccessRefreshToken(user)
	if err != nil {
		log.Println("Ошибка при генерации токенов:", err)
		return nil, err
	}

	if _, err := users.UpdateUserRefreshToken(manager, user.Guid, refreshTokenHash); err != nil {
		log.Println("Ошибка при обновлении пользователя:", err)
		return nil, err
	}

	responseData := map[string]interface{}{
		"accessToken":  accesTokenString,
		"refreshToken": refreshToken,
	}

	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		log.Println("Ошибка при сериализации ответа:", err)
		return nil, err
	}

	return jsonResponse, nil
}
