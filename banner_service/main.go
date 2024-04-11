// В файле main.go

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/15Andrew43/backend-trainee-assignment-2024/config"
	"github.com/15Andrew43/backend-trainee-assignment-2024/database"
	"github.com/15Andrew43/backend-trainee-assignment-2024/handlers"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	if err := database.ConnectToPostgres(&config.Cfg); err != nil {
		log.Fatalf("Ошибка при кодключении к Postgres: %v", err)
	}
	defer database.PgConn.Close(context.Background())

	if err := database.ConnectToMongoDB(&config.Cfg); err != nil {
		log.Fatalf("Ошибка при кодключении к Mongo: %v", err)
	}
	defer database.MongoCli.Disconnect(context.Background())

	http.HandleFunc("/", handlers.GetUserBanner)

	log.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
