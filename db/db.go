package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Manager struct {
	Conn *sql.DB
}

func NewDBManager(driverName string, sourceName string) *Manager {
	var db *sql.DB
	var err error

	const maxRetries = 3
	const retryDelay = 2 * time.Second

	for i := 1; i <= maxRetries; i++ {
		db, err = sql.Open(driverName, sourceName)
		if err != nil {
			log.Printf("Попытка %d: ошибка при открытии соединения: %v", i, err)
		} else if err = db.Ping(); err != nil {
			log.Printf("Попытка %d: база данных не отвечает: %v", i, err)
		} else {
			log.Println("Успешное подключение к базе данных")
			return &Manager{Conn: db}
		}

		time.Sleep(retryDelay)
	}

	panic(fmt.Sprintf("Не удалось подключиться к базе данных после %d попыток: %v", maxRetries, err))
}

func (manager *Manager) Close() {
	if manager.Conn != nil {
		manager.Conn.Close()
		manager.Conn = nil
	}
}