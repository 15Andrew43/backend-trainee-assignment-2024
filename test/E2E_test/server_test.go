package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/15Andrew43/backend-trainee-assignment-2024/client"
	"github.com/15Andrew43/backend-trainee-assignment-2024/fakes"
	"github.com/15Andrew43/backend-trainee-assignment-2024/model"
)

func TestE2E(t *testing.T) {
	fmt.Println("Starting test")

	baseURL := "http://localhost:8080"
	adminToken := fakes.Admin
	userToken := fakes.AuthorizedUser

	// 1. Попытка добавить баннер от лица не админа
	resp, err := client.AddBanner(baseURL, userToken, model.BannerRequest{})
	if err != nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Ожидается статус код %d, получено %d", http.StatusForbidden, resp.StatusCode)
	}
	resp.Body.Close()
	fmt.Println("1. Попытка добавить баннер от лица не админа успешно обработана")

	// // 2. Добавление двух баннеров с пересекающимися тэгами и одного баннера с другими тэгами от админа
	banner1 := model.BannerRequest{TagIds: []int{1}, FeatureId: 1, Content: "{\"title\": \"banner1\"}"}
	banner2 := model.BannerRequest{TagIds: []int{1, 2}, FeatureId: 2, Content: "{\"title\": \"banner2\"}"}
	banner3 := model.BannerRequest{TagIds: []int{3}, FeatureId: 3, Content: "{\"title\": \"banner3\"}"}

	resp, err = client.AddBanner(baseURL, adminToken, banner1)
	if err != nil {
		t.Fatalf("Ошибка при добавлении баннера: %v", err)
	}
	resp.Body.Close()
	resp, err = client.AddBanner(baseURL, adminToken, banner2)
	if err != nil {
		t.Fatalf("Ошибка при добавлении баннера: %v", err)
	}
	resp.Body.Close()
	resp, err = client.AddBanner(baseURL, adminToken, banner3)
	if err != nil {
		t.Fatalf("Ошибка при добавлении баннера: %v", err)
	}
	resp.Body.Close()
	fmt.Println("2. Добавление баннеров завершено")

	// 3. Попытка получить баннер от имени пользователя с несовпадающим тегом+фичей
	resp, err = client.GetUserBanner(baseURL, userToken, 2, 1, true)
	if err != nil || resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Ожидается статус код %d, получено %d", http.StatusNotFound, resp.StatusCode)
	}
	resp.Body.Close()
	fmt.Println("3. Попытка получить баннер от пользователя с несовпадающим тегом успешно обработана")

	// 4. Получение баннера от имени пользователя с совпадающим тегом
	resp, err = client.GetUserBanner(baseURL, userToken, 1, 1, true)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидается статус код %d, получено %d ", http.StatusOK, resp.StatusCode)
	}
	resp.Body.Close()
	fmt.Println("4. Получение баннера от пользователя с совпадающим тегом успешно завершено")

	// 5. Получение баннера от имени админа
	resp, err = client.GetUserBanner(baseURL, userToken, 1, 1, true)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидается статус код %d, получено %d ", http.StatusOK, resp.StatusCode)
	}
	resp.Body.Close()
	fmt.Println("5. Получение баннера от пользователя с совпадающим тегом успешно завершено")

	// 6. Попытка изменить баннер от пользователя без подходящего тега
	resp, err = client.UpdateBanner(baseURL, userToken, 1, model.BannerRequest{TagIds: []int{1}, FeatureId: 3, Content: "{\"title\": \"new feature 3 in banner\"}"})
	if err != nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Ожидается статус код %d, получено %d", http.StatusForbidden, resp.StatusCode)
	}
	resp.Body.Close()
	fmt.Println("6. Попытка изменить баннер от пользователя без подходящего тега успешно обработана")

	// 7. Изменение баннера от админа
	resp, err = client.UpdateBanner(baseURL, adminToken, 1, model.BannerRequest{TagIds: []int{1}, FeatureId: 3, Content: "{\"title\": \"new feature 3 in banner\"}"})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидается статус код %d, получено %d ", http.StatusOK, resp.StatusCode)
	}
	resp.Body.Close()
	fmt.Println("7. Изменение баннера от пользователя с подходящим тегом успешно завершено")

	// time.Sleep(6 * time.Second) // for checking that redis works ok

	// 8. Проверка изменений баннера
	resp, err = client.GetUserBanner(baseURL, userToken, 2, 2, true)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидается статус код %d, получено %d ", http.StatusOK, resp.StatusCode)
	}
	var updatedBanner model.MongoBannerData

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Ошибка при чтении тела ответа: %v", err)
	}

	if err := json.Unmarshal(body, &updatedBanner); err != nil {
		t.Fatalf("Ошибка при десериализации баннера из тела ответа: %v", err)
	}

	fmt.Printf("updatedBanner = %+v\n\n", updatedBanner)

	var contentExpected map[string]interface{}

	err = json.Unmarshal([]byte(banner2.Content), &contentExpected)
	if err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return
	}

	if !reflect.DeepEqual(updatedBanner.Content, contentExpected) {
		t.Fatalf("Содержимое обновленного баннера не соответствует ожидаемому")
	}
	resp.Body.Close()
	fmt.Println("8. Проверка изменений баннера успешно завершена")

	resp, err = client.GetBanners(baseURL, userToken, 1, 1)
	if err != nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Ожидается ошибка доступа при выполнении GetBanners от лица пользователя, но ошибка не возникла")
	}
	resp.Body.Close()

	// 9. Проверка изменений баннера
	resp, err = client.GetBanners(baseURL, adminToken, 3, 3)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидается статус код %d, получено %d ", http.StatusOK, resp.StatusCode)
	}
	var updatedBanners []model.MongoBannerData

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Ошибка при чтении тела ответа: %v", err)
	}

	if err := json.Unmarshal(body, &updatedBanners); err != nil {
		t.Fatalf("Ошибка при десериализации баннеров из тела ответа: %v", err)
	}

	if len(updatedBanners) > 0 {
		var contentExpected map[string]interface{}

		err = json.Unmarshal([]byte(banner3.Content), &contentExpected)
		if err != nil {
			t.Fatalf("Ошибка при разборе JSON из ожидаемого баннера: %v", err)
		}

		if !reflect.DeepEqual(updatedBanners[0].Content, contentExpected) {
			t.Fatalf("Содержимое обновленного баннера не соответствует ожидаемому")
		}
	} else {
		t.Fatalf("Не найдено обновленных баннеров")
	}
	resp.Body.Close()
	fmt.Println("9. Проверка изменений баннера успешно завершена")

	// 10. Удаление баннера от пользователя с несовпадающим тегом
	resp, err = client.DeleteBanner(baseURL, userToken, 1)
	if err != nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Ожидается статус код %d, получено %d", http.StatusForbidden, resp.StatusCode)
	}
	fmt.Println("10. Удаление баннера от пользователя с несовпадающим тегом успешно обработано")

	// 11. Удаление оставшихся двух баннеров от админа
	resp, err = client.DeleteBanner(baseURL, adminToken, 1)
	if err != nil || resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Ожидается статус код %d, получено %d ", http.StatusNoContent, resp.StatusCode)
	}
	resp, err = client.DeleteBanner(baseURL, adminToken, 2)
	if err != nil || resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Ожидается статус код %d, получено %d ", http.StatusNoContent, resp.StatusCode)
	}
	resp, err = client.DeleteBanner(baseURL, adminToken, 3)
	if err != nil || resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Ожидается статус код %d, получено %d ", http.StatusNoContent, resp.StatusCode)
	}
	fmt.Println("11. Удаление оставшихся двух баннеров от админа успешно завершено")

	// time.Sleep(6 * time.Second) // for checking that redis works ok

	// 12. Проверка отсутствия оставшихся баннеров
	resp, err = client.GetUserBanner(baseURL, adminToken, 1, 1, true)
	if err != nil || resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Ожидается статус код %d, получено %d", http.StatusNotFound, resp.StatusCode)
	}
	fmt.Println("12. Проверка отсутствия оставшихся баннеров успешно завершена")
}
