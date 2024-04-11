package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/15Andrew43/backend-trainee-assignment-2024/database"
	"github.com/15Andrew43/backend-trainee-assignment-2024/model"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserBanner(w http.ResponseWriter, r *http.Request) {
	// postgres
	tagID, _ := strconv.Atoi(r.URL.Query().Get("tag_id"))
	featureID, _ := strconv.Atoi(r.URL.Query().Get("feature_id"))

	var banner model.Banner
	err := database.GetPostgresBanner(tagID, featureID, &banner)
	log.Printf("type err = %T", err)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Printf("Не найдено строк с tag_id = %d и feature_id = %d", tagID, featureID)
			http.Error(w, "Баннер не найден в Postgres", http.StatusNotFound)
			return
		}
		log.Printf("Ошибка при выполнении запроса к Postgres: %v", err)
		http.Error(w, "Внутренняя ошибка сервера при запросе к Postgres", http.StatusInternalServerError)
		return
	}

	// mongo
	var bannerData model.BannerData
	err = database.GetMongoBannerData(&bannerData, &banner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Не найдено документов с data_id = %v", banner.DataID)
			http.Error(w, "Баннер не найден в Mongo", http.StatusNotFound)
			return
		}
		log.Printf("Ошибка при выполнении запроса к MongoDB: %v", err)
		http.Error(w, "Внутренняя ошибка сервера при запросе к Mongo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bannerData)
}