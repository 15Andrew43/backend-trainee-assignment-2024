package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/15Andrew43/backend-trainee-assignment-2024/database"
	"github.com/15Andrew43/backend-trainee-assignment-2024/model"
	myerrors "github.com/15Andrew43/backend-trainee-assignment-2024/my_errors"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserBanner(w http.ResponseWriter, r *http.Request) {
	// postgres
	tagID, _ := strconv.Atoi(r.URL.Query().Get("tag_id"))
	featureID, _ := strconv.Atoi(r.URL.Query().Get("feature_id"))

	var banner model.Banner
	err := database.GetPostgresBanner(tagID, featureID, &banner)
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

func GetAllBanners(w http.ResponseWriter, r *http.Request) {
	// postgres
	tagID := -1
	if tagStr := r.URL.Query().Get("tag_id"); tagStr != "" {
		var error error
		tagID, error = strconv.Atoi(tagStr)
		if error != nil {
			log.Printf("Некорректные данные tagID: %v", error)
			http.Error(w, "Некорректные данные", http.StatusBadRequest)
			return
		}
	}

	featureID := -1
	if featureStr := r.URL.Query().Get("feature_id"); featureStr != "" {
		var error error
		featureID, error = strconv.Atoi(featureStr)
		if error != nil {
			log.Printf("Некорректные данные featureID: %v", error)
			http.Error(w, "Некорректные данные", http.StatusBadRequest)
			return
		}
	}

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		var error error
		limit, error = strconv.Atoi(limitStr)
		if error != nil {
			log.Printf("Некорректные данные limit: %v", error)
			http.Error(w, "Некорректные данные", http.StatusBadRequest)
			return
		}
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		var error error
		offset, error = strconv.Atoi(offsetStr)
		if error != nil {
			log.Printf("Некорректные данные offset: %v", error)
			http.Error(w, "Некорректные данные", http.StatusBadRequest)
			return
		}
	}

	banners, err := database.GetPostgresAllBanners(tagID, featureID, limit, offset)
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
	var bannerDatas []model.BannerData
	for _, banner := range banners {
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
		bannerDatas = append(bannerDatas, bannerData)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bannerDatas)
}

func CreateBanner(w http.ResponseWriter, r *http.Request) {

	requestBody, ok := r.Context().Value("requestBody").(model.RequestBodyBanner)
	if !ok {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	nextId, err := database.CreatePostgresBanner(&requestBody)
	if err != nil {
		var bannerExistErr *myerrors.BannerExist
		if errors.As(err, &bannerExistErr) {
			log.Printf("Ошибка при вставке данных в Postgres: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Ошибка при вставке данных в Postgres: %v", err)
		http.Error(w, "Внутренняя ошибка сервера при запросе к Postgres", http.StatusInternalServerError)
		return
	}

	err = database.CreateMongoBanner(nextId, requestBody.Content)
	if err != nil {
		log.Printf("Ошибка при вставке данных в Mongo: %v", err)
		http.Error(w, "Внутренняя ошибка сервера при запросе к Mongo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Ошибка при ковертации строки %s в число", vars["id"])
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestBody, ok := r.Context().Value("requestBody").(model.RequestBodyBanner)
	if !ok {
		log.Printf("heer is DEBUGGING")
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	dataId, err := database.UpgradePostgresBanner(id, &requestBody)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Printf("Не найдено строк с banner_id = %d", id)
			http.Error(w, "Баннер не найден в Postgres при обновлении", http.StatusNotFound)
			return
		}
		log.Printf("Ошибка при обновлении данных в Postgres: %v", err)
		http.Error(w, "Внутренняя ошибка сервера при запросе к Postgres", http.StatusInternalServerError)
		return
	}

	err = database.UpgradeMongoBanner(dataId, requestBody.Content)
	if err != nil {
		log.Printf("Ошибка при обновлении данных в Mongo: %v", err)
		http.Error(w, "Внутренняя ошибка сервера при запросе к Mongo", http.StatusInternalServerError)
		return
	}
}

func DeleteBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Ошибка при ковертации строки %s в число", vars["id"])
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataId, err := database.DeletePostgresBanner(id)
	if err != nil {
		log.Printf("Ошибка при обновлении данных в Postgres: %v", err)
		http.Error(w, "Внутренняя ошибка сервера при запросе к Postgres", http.StatusInternalServerError)
		return
	}

	err = database.DeleteMongoBanner(dataId)
	if err != nil {
		log.Printf("Ошибка при обновлении данных в Mongo: %v", err)
		http.Error(w, "Внутренняя ошибка сервера при запросе к Mongo", http.StatusInternalServerError)
		return
	}
}
