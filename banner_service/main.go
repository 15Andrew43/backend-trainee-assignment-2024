package main

import (
	"context"
	"log"
	"net/http"

	"github.com/15Andrew43/backend-trainee-assignment-2024/config"
	"github.com/15Andrew43/backend-trainee-assignment-2024/database"
	"github.com/15Andrew43/backend-trainee-assignment-2024/handlers"
	"github.com/gorilla/mux"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	if err := database.ConnectToPostgres(&config.Cfg); err != nil {
		log.Fatalf("Ошибка при кодключении к Postgres: %v", err)
	}
	defer database.PgPool.Close()

	if err := database.ConnectToMongoDB(&config.Cfg); err != nil {
		log.Fatalf("Ошибка при кодключении к Mongo: %v", err)
	}
	defer database.MongoCli.Disconnect(context.Background())

	r := mux.NewRouter()

	r.Handle("/user_banner", handlers.UserBannerHandler).Methods("GET")
	r.Handle("/banner", handlers.BannerHandler).Methods("GET")
	r.Handle("/banner", handlers.CreateBannerHandler).Methods("POST")
	r.Handle("/banner/{id}", handlers.UpdateBannerHandler).Methods("PATCH")
	r.Handle("/banner/{id}", handlers.DeleteBannerHandler).Methods("DELETE")

	log.Println("Server is listening on port 8080...")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Ошибка во время запуска сервера: ", err)
	}
}
