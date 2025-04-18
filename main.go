package main

import (
	"log"
	"med/services"
	"med/db"
	"os"
)

func main() {
	services.LoadEnvFile(".env")
	manager := db.NewDBManager("postgres", os.Getenv("DB_CONNECTION_STRING"))
	db.Migrate(manager)

	log.Println("База данных успешно инициализирована и мигрирована.")
	defer manager.Close()
}