package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"med/db"
	"med/db/tokens"
	"med/models"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var secret = []byte("super-duper-ultra-mega-secret-pentagon-key")

func GenerateAccessRefreshToken(user models.User) (string, string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id":   user.Guid,
		"ip":        user.Ip,
		"access_id": user.AccessId,
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
	})

	accessToken, err := token.SignedString(secret)
	if err != nil {
		return "", "", "", err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return "", "", "", err
	}
	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, string(refreshTokenHash), nil
}

func UpdateTokens(manager *db.Manager, user models.User, refreshToken string, accessToken string) (string, string, string, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("неправильный алгоритм подписи")
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return "", "", "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)
	accessID := claims["access_id"].(string)
	tokenIP := claims["ip"].(string)

	if tokenIP != user.Ip {
		log.Println("Не совпадают IP адреса, отправляю письмо на почту пользователя...", user.Email)
	}

	if accessID != user.AccessId {
		return "", "", "", fmt.Errorf("неправильный ключ доступа")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.RefreshToken), []byte(refreshToken))
	if err != nil {
		return "", "", "", fmt.Errorf("неправильный ключ обновления")
	}

	if user.Guid != userID {
		return "", "", "", fmt.Errorf("неправильный идентификатор пользователя")
	}

	isUsed, err := tokens.IsTokenUsed(manager, user.Guid, user.RefreshToken)
	if err != nil {
		return "", "", "", fmt.Errorf("ошибка при проверке использованного токена: %w", err)
	}
	if isUsed {
		return "", "", "", fmt.Errorf("refresh token уже был использован")
	}

	if err := tokens.MarkTokenAsUsed(manager, user.Guid, user.RefreshToken); err != nil {
		return "", "", "", fmt.Errorf("ошибка при маркировке токена как использованного: %w", err)
	}

	return GenerateAccessRefreshToken(user)
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
