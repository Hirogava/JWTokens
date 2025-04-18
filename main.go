package main

import (
	"fmt"
	"log"
	"med/db"
	"med/db/tokens"
	"med/routes"
	"med/services"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	services.LoadEnvFile(".env")
	manager := db.NewDBManager("postgres", os.Getenv("DB_CONNECTION_STRING"))
	db.Migrate(manager)

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := tokens.CleanupOldTokens(manager); err != nil {
					log.Printf("Ошибка при очистке старых токенов: %v\n", err)
				}
			}
		}
	}()

	log.Println("База данных успешно инициализирована и мигрирована.")
	defer manager.Close()

	r := mux.NewRouter()

	routes.Init(r, manager)

	serverPort := os.Getenv("SERVER_PORT")
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: r,
	}

	log.Println("Сервер запущен на порту " + serverPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
