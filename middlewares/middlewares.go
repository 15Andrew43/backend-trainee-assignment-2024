package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/15Andrew43/backend-trainee-assignment-2024/fakes"
	"github.com/15Andrew43/backend-trainee-assignment-2024/model"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	Query  []string
	Header []string
	Data   []string
}

func CheckParamsMiddleware(needParams Params) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, param := range needParams.Query {
				if r.URL.Query().Get(param) == "" {
					http.Error(w, fmt.Sprintf("Отсутствует обязательный параметр %v в параметрах пути", param), http.StatusBadRequest)
					return
				}
			}
			for _, param := range needParams.Header {
				if r.Header.Get(param) == "" {
					http.Error(w, fmt.Sprintf("Отсутствует обязательный параметр %v в параметрах пути", param), http.StatusBadRequest)
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
				log.Printf("body = %+v\n\n\n\n\n", body)
				log.Printf("type of featureId = %T\n\n", body["feature_id"])
				log.Printf("type of tags = %T\n\n", body["tag_ids"])

				var requestBody model.RequestBodyBanner
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
						content, ok := value.(string)
						if !ok {
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

				log.Printf("body = %+v\n\n\n\n\n", requestBody)

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
