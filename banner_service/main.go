package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/15Andrew43/backend-trainee-assignment-2024/config"
	"github.com/15Andrew43/backend-trainee-assignment-2024/database"
	"github.com/15Andrew43/backend-trainee-assignment-2024/handlers"
	"github.com/15Andrew43/backend-trainee-assignment-2024/middlewares"
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
	defer database.PgConn.Close(context.Background())

	if err := database.ConnectToMongoDB(&config.Cfg); err != nil {
		log.Fatalf("Ошибка при кодключении к Mongo: %v", err)
	}
	defer database.MongoCli.Disconnect(context.Background())

	r := mux.NewRouter()

	r.HandleFunc("/latest/{id}", func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		fmt.Println("here ", id)
	})

	userBannerHandler := middlewares.CheckParamsMiddleware(middlewares.Params{
		Query:  []string{"tag_id", "feature_id"},
		Header: []string{"token"},
	})(http.HandlerFunc(handlers.GetUserBanner))

	bannerHandler := middlewares.CheckParamsMiddleware(middlewares.Params{
		Header: []string{"token"},
	})(http.HandlerFunc(handlers.GetAllBanners))

	createBannerHandler := middlewares.CheckParamsMiddleware(middlewares.Params{
		Header: []string{"token"},
		Data:   []string{"tag_ids", "feature_id", "content", "is_active"},
	})(http.HandlerFunc(handlers.CreateBanner))

	updateBannerHandler := middlewares.CheckParamsMiddleware(middlewares.Params{
		URLParam: []string{"id"},
		Header:   []string{"token"},
		Data:     []string{"tag_ids", "feature_id", "content", "is_active"},
	})(http.HandlerFunc(handlers.UpdateBanner))

	deleteBannerHandler := middlewares.CheckParamsMiddleware(middlewares.Params{
		URLParam: []string{"id"},
		Header:   []string{"token"},
	})(http.HandlerFunc(handlers.DeleteBanner))

	r.Handle("/user_banner", userBannerHandler).Methods("GET")
	r.Handle("/banner", bannerHandler).Methods("GET")
	r.Handle("/banner", createBannerHandler).Methods("POST")
	r.Handle("/banner/{id}", updateBannerHandler).Methods("PATCH")
	r.Handle("/banner/{id}", deleteBannerHandler).Methods("DELETE")

	log.Println("Server is listening on port 8080...")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Ошибка во время запуска сервера: ", err)
	}
}
