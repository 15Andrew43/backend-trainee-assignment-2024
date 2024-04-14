package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/15Andrew43/backend-trainee-assignment-2024/database"
	"github.com/15Andrew43/backend-trainee-assignment-2024/middlewares"
	"github.com/15Andrew43/backend-trainee-assignment-2024/model"
	myerrors "github.com/15Andrew43/backend-trainee-assignment-2024/my_errors"
	"github.com/15Andrew43/backend-trainee-assignment-2024/util"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
)

var UserBannerHandler = middlewares.AuthMiddleware(middlewares.CheckParamsMiddleware(middlewares.Params{
	Query:  []string{"tag_id", "feature_id"},
	Header: []string{"token"},
})(http.HandlerFunc(GetUserBanner)))

var BannerHandler = middlewares.AuthMiddleware(middlewares.CheckParamsMiddleware(middlewares.Params{
	Header: []string{"token"},
})(http.HandlerFunc(GetAllBanners)))

var CreateBannerHandler = middlewares.AuthMiddleware(middlewares.CheckParamsMiddleware(middlewares.Params{
	Header: []string{"token"},
	Data:   []string{"tag_ids", "feature_id", "content", "is_active"},
})(http.HandlerFunc(CreateBanner)))

var UpdateBannerHandler = middlewares.AuthMiddleware(middlewares.CheckParamsMiddleware(middlewares.Params{
	URLParam: []string{"id"},
	Header:   []string{"token"},
	Data:     []string{"tag_ids", "feature_id", "content", "is_active"},
})(http.HandlerFunc(UpdateBanner)))

var DeleteBannerHandler = middlewares.AuthMiddleware(middlewares.CheckParamsMiddleware(middlewares.Params{
	URLParam: []string{"id"},
	Header:   []string{"token"},
})(http.HandlerFunc(DeleteBanner)))

func GetUserBanner(w http.ResponseWriter, r *http.Request) {
	// postgres
	tagID, _ := strconv.Atoi(r.URL.Query().Get("tag_id"))
	featureID, _ := strconv.Atoi(r.URL.Query().Get("feature_id"))

	errorPostgresChan := make(chan error, 1)

	errorMongoChan := make(chan myerrors.DataError, 1)

	wg := sync.WaitGroup{}

	var banner model.PostgresBanner
	var bannerData model.MongoBannerData
	wg.Add(1)
	go func() {
		defer wg.Done()
		database.GetPostgresBanner(tagID, featureID, errorPostgresChan, &banner)
	}()

	err := <-errorPostgresChan
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

	wg.Add(1)
	go func() {
		wg.Done()
		database.GetMongoBannerData(errorMongoChan, &bannerData, &banner)
	}()

	wg.Wait()

	var dataErr myerrors.DataError = <-errorMongoChan
	if dataErr.Err != nil {
		if dataErr.Err == mongo.ErrNoDocuments {
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
	log.Printf("GetPostgresAllBanners is done???")
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

	var errorMongoChans []chan myerrors.DataError = make([]chan myerrors.DataError, len(banners))
	for i := 0; i < len(banners); i++ {
		errorMongoChans[i] = make(chan myerrors.DataError, 1)
	}

	wg := sync.WaitGroup{}

	semaphore := make(chan struct{}, 20)

	var bannerDatas []model.MongoBannerData = make([]model.MongoBannerData, len(banners))
	for i, banner := range banners {
		semaphore <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				<-semaphore
			}()
			database.GetMongoBannerData(errorMongoChans[i], &bannerDatas[i], &banner)
		}()
	}

	wg.Wait()

	for _, ch := range errorMongoChans {
		var err myerrors.DataError = <-ch
		if err.Err != nil {
			if err.Err == mongo.ErrNoDocuments {
				log.Printf("Не найдено документов с data_id = %v", err.DataID)
				http.Error(w, "Баннер не найден в Mongo", http.StatusNotFound)
				return
			}
			log.Printf("Ошибка при выполнении запроса к MongoDB: %v", err)
			http.Error(w, "Внутренняя ошибка сервера при запросе к Mongo", http.StatusInternalServerError)
			return
		}
	}

	if len(bannerDatas) == 0 {
		http.Error(w, "Таких баннеров не найдено", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bannerDatas)
}

func CreateBanner(w http.ResponseWriter, r *http.Request) {

	requestBody, ok := r.Context().Value("requestBody").(model.Banner)
	if !ok {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	nextId := util.GenerateNextId()

	errorPostgresChan := make(chan error, 1)

	errorMongoChan := make(chan error, 1)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		database.CreatePostgresBanner(nextId, errorPostgresChan, &requestBody)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		database.CreateMongoBanner(nextId, errorMongoChan, requestBody.Content)
	}()

	wg.Wait()

	err := <-errorPostgresChan
	/////////////////POSTGRES//////////////////
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

	err = <-errorMongoChan
	////////////////MONGO///////////
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

	requestBody, ok := r.Context().Value("requestBody").(model.Banner)
	if !ok {
		log.Printf("heer is DEBUGGING")
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	errorPostgresChan := make(chan error, 1)

	chanDataId := make(chan int, 1)

	errorMongoChan := make(chan error, 1)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		database.UpgradePostgresBanner(id, errorPostgresChan, chanDataId, &requestBody)
	}()

	dataId := <-chanDataId

	wg.Add(1)
	go func() {
		defer wg.Done()
		database.UpgradeMongoBanner(dataId, errorMongoChan, requestBody.Content)
	}()

	wg.Wait()

	err = <-errorPostgresChan
	if err != nil {
		var bannerExistErr *myerrors.BannerExist
		if errors.As(err, &bannerExistErr) {
			log.Printf("Ошибка при вставке данных в Postgres: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Printf("Не найдено строк с banner_id = %d", id)
			http.Error(w, "Баннер не найден в Postgres при обновлении", http.StatusNotFound)
			return
		}
		log.Printf("Ошибка при обновлении данных в Postgres: %v", err)
		http.Error(w, "Внутренняя ошибка сервера при запросе к Postgres", http.StatusInternalServerError)
		return
	}

	err = <-errorMongoChan
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

	errorPostgresChan := make(chan myerrors.DataError, 1)

	errorMongoChan := make(chan error, 1)

	chanDataId := make(chan int, 1)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		database.DeletePostgresBanner(id, chanDataId, errorPostgresChan)
	}()

	dataId := <-chanDataId

	wg.Add(1)
	go func() {
		defer wg.Done()
		database.DeleteMongoBanner(dataId, errorMongoChan)
	}()

	wg.Wait()

	var dataErr = <-errorPostgresChan
	if dataErr.Err != nil {
		if strings.Contains(dataErr.Err.Error(), "no rows in result set") {
			log.Printf("Не найдено строк с id = %d", dataErr.DataID)
			http.Error(w, "Строки с таким id не было", http.StatusBadRequest)
			return
		}
		log.Printf("Ошибка при обновлении данных в Postgres: %v", err)
		http.Error(w, "Внутренняя ошибка сервера при запросе к Postgres", http.StatusInternalServerError)
		return
	}
	err = <-errorMongoChan
	if err != nil {
		log.Printf("Ошибка при обновлении данных в Mongo: %v", err)
		http.Error(w, "Внутренняя ошибка сервера при запросе к Mongo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
