package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/15Andrew43/backend-trainee-assignment-2024/database"
	"github.com/15Andrew43/backend-trainee-assignment-2024/fakes"
	"github.com/15Andrew43/backend-trainee-assignment-2024/model"
	"github.com/gorilla/mux"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("--------------------- id = %v", chi.URLParam(r, "id"))
		token := r.Header.Get("token")

		if token == "" {
			http.Error(w, "Не указан заголовок token", http.StatusUnauthorized)
			return
		}
		role := fakes.GetRole(token)
		if role == fakes.NotAuthorizedUser || (role == fakes.AuthorizedUser && r.URL.Path != "/user_banner") {
			http.Error(w, "Пользователь не имеет доступа", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type Params struct {
	Query    []string
	URLParam []string
	Header   []string
	Data     []string
}

func CheckParamsMiddleware(needParams Params) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// log.Printf("path id = %v", chi.URLParam(r, "id")) // всегда выводуит пустую строку! :(
			for _, param := range needParams.Query {
				if r.URL.Query().Get(param) == "" {
					http.Error(w, fmt.Sprintf("Отсутствует обязательный параметр %v в параметрах запроса", param), http.StatusBadRequest)
					return
				}
			}
			for _, param := range needParams.URLParam {
				vars := mux.Vars(r)
				if _, ok := vars[param]; !ok {
					http.Error(w, fmt.Sprintf("Отсутствует обязательный параметр %v в параметрах пути", param), http.StatusBadRequest)
					return
				}
			}
			for _, param := range needParams.Header {
				if r.Header.Get(param) == "" {
					http.Error(w, fmt.Sprintf("Отсутствует обязательный параметр %v в заголовках", param), http.StatusBadRequest)
					return
				}
			}
			if len(needParams.Data) > 0 {
				body := map[string]interface{}{}
				err := json.NewDecoder(r.Body).Decode(&body)
				if err != nil {
					http.Error(w, "Ошибка при чтении JSON-тела запроса", http.StatusBadRequest)
					return
				}

				var requestBody model.Banner
				for key, value := range body {
					switch key {
					case "tag_ids":
						var tagIDs []int
						for _, tagID := range value.([]interface{}) {
							tag, ok := tagID.(float64)
							if !ok {
								http.Error(w, "Некорректные данные для tag_ids", http.StatusBadRequest)
								return
							}
							tagIDs = append(tagIDs, int(tag))
						}
						requestBody.TagIds = tagIDs
					case "feature_id":
						featureID, ok := value.(float64)
						if !ok {
							http.Error(w, "Некорректные данные для feature_id", http.StatusBadRequest)
							return
						}
						requestBody.FeatureId = int(featureID)
					case "content":
						contentStr, ok := value.(string)
						if !ok {
							log.Printf("Ошибка во время value.(string)")
							http.Error(w, "Некорректные данные для content", http.StatusBadRequest)
							return
						}

						var content map[string]interface{}
						err := json.Unmarshal([]byte(contentStr), &content)
						if err != nil {
							log.Printf("Ошибка во время Unmarshal: %v", err)
							http.Error(w, "Некорректные данные для content", http.StatusBadRequest)
							return
						}
						requestBody.Content = content
					case "is_active":
						isActive, ok := value.(bool)
						if !ok {
							http.Error(w, "Некорректные данные для is_active", http.StatusBadRequest)
							return
						}
						requestBody.IsActive = isActive
					}
				}

				for _, param := range needParams.Data {
					_, ok := body[param]
					if !ok {
						http.Error(w, fmt.Sprintf("Отсутствует обязательный параметр %v в теле запроса", param), http.StatusBadRequest)
						return
					}
				}
				r = r.WithContext(context.WithValue(r.Context(), "requestBody", requestBody))
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		useLastRevision := false
		useLastRevisionStr := r.URL.Query().Get("use_last_revision")
		if useLastRevisionStr != "" {
			var err error
			useLastRevision, err = strconv.ParseBool(useLastRevisionStr)
			if err != nil {
				http.Error(w, "Некорректные данные для useLastRevision", http.StatusBadRequest)
				return
			}
		}

		key := r.URL.Path

		if !useLastRevision {
			val, err := database.RedisClient.Get(context.Background(), key).Result()
			if err == nil {
				log.Println("Данные найдены в кеше")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(val))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func SaveToCache(key string, data interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = database.RedisClient.Set(context.Background(), key, jsonData, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
