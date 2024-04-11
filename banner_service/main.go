package main

import (
	"context"
	"log"
	"net/http"

	"github.com/15Andrew43/backend-trainee-assignment-2024/config"
	"github.com/15Andrew43/backend-trainee-assignment-2024/database"
	"github.com/15Andrew43/backend-trainee-assignment-2024/handlers"
	"github.com/15Andrew43/backend-trainee-assignment-2024/middlewares"
	"github.com/go-chi/chi/v5"
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

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)

		r.With(middlewares.CheckParamsMiddleware(middlewares.Params{Query: []string{"tag_id", "feature_id"}, Header: []string{"token"}})).Get("/user_banner", handlers.GetUserBanner)

		r.With(middlewares.CheckParamsMiddleware(middlewares.Params{Header: []string{"token"}})).Get("/banner", handlers.GetAllBanners)

		r.With(middlewares.CheckParamsMiddleware(middlewares.Params{Header: []string{"token"}, Data: []string{"tag_ids", "feature_id", "content", "is_active"}})).Post("/banner", handlers.CreateBanner)
	})

	log.Println("Server is listening on port 8080...")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Ошибка во время запуска сервера: %v", err)
	}
}
