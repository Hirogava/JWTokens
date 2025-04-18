package users

import (
	"med/db"
	"med/models"
)

func GetUserByGuid(guid string, manager *db.Manager) (models.User, error) {
	var user models.User

	query := `SELECT * FROM users WHERE guid=$1`
	
	err := manager.Conn.QueryRow(query, guid).Scan(&user.Guid, &user.Ip, &user.Email, &user.RefreshToken, &user.AccessId, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func CreateUser(manager *db.Manager, ip string, email string) (models.User, error) {
	var user models.User

	query := `INSERT INTO users(ip, email) VALUES($1, $2) RETURNING guid, ip, email, access_token, created_at`
	err := manager.Conn.QueryRow(query, ip, email).Scan(&user.Guid, &user.Ip, &user.Email, &user.AccessId, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}
	
	return user, nil
}

func UpdateUserRefreshToken(manager *db.Manager, guid string, refreshToken string) (models.User, error) {
	var user models.User

	query := `UPDATE users SET refresh_token=$1 WHERE guid=$2 RETURNING *`
	err := manager.Conn.QueryRow(query, refreshToken, guid).Scan(&user.Guid, &user.Ip, &user.Email, &user.RefreshToken, &user.AccessId, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}