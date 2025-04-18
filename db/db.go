package db

import (
	"database/sql"
	"fmt"
)

type Manager struct {
	Conn *sql.DB
}

func NewDBManager(driverName string, sourceName string) *Manager {
	db, err := sql.Open(driverName, sourceName)
	if err != nil {
		panic(fmt.Sprintf("Не удалось подключиться к базе данных: %v", err))
	}

	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("База данных не отвечает: %v", err))
	}

	return &Manager{
		Conn: db,
	}
}

func (manager *Manager) Close() {
	if manager.Conn != nil {
		manager.Conn.Close()
		manager.Conn = nil
	}
}