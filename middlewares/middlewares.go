package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/15Andrew43/backend-trainee-assignment-2024/fakes"
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
				var _, bodyParams map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&bodyParams)
				if err != nil {
					http.Error(w, "Ошибка при чтении JSON-тела запроса", http.StatusBadRequest)
					return
				}

				for _, param := range needParams.Data {
					if _, ok := bodyParams[param]; !ok {
						http.Error(w, fmt.Sprintf("Отсутствует обязательный параметр %v в параметрах пути", param), http.StatusBadRequest)
						return
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
