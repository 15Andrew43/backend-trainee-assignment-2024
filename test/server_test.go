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

	// // 1. Попытка добавить баннер от лица не админа
	// resp, err := client.AddBanner(baseURL, userToken, model.BannerRequest{})
	// if resp.StatusCode != http.StatusForbidden {
	// 	t.Fatalf("Ожидается статус код %d, получено %d", http.StatusForbidden, resp.StatusCode)
	// }
	// fmt.Println("1. Попытка добавить баннер от лица не админа успешно обработана")

	// // 2. Добавление двух баннеров с пересекающимися тэгами и одного баннера с другими тэгами от админа
	banner1 := model.BannerRequest{TagIds: []int{1}, FeatureId: 1, Content: "{\"title\": \"banner1\"}"}
	banner2 := model.BannerRequest{TagIds: []int{1, 2}, FeatureId: 2, Content: "{\"title\": \"banner2\"}"}
	banner3 := model.BannerRequest{TagIds: []int{3}, FeatureId: 3, Content: "{\"title\": \"banner3\"}"}
	_ = banner1
	_ = banner3

	// resp, err = client.AddBanner(baseURL, adminToken, banner1)
	// if err != nil {
	// 	t.Fatalf("Ошибка при добавлении баннера: %v", err)
	// }
	// resp, err = client.AddBanner(baseURL, adminToken, banner2)
	// if err != nil {
	// 	t.Fatalf("Ошибка при добавлении баннера: %v", err)
	// }
	// resp, err = client.AddBanner(baseURL, adminToken, banner3)
	// if err != nil {
	// 	t.Fatalf("Ошибка при добавлении баннера: %v", err)
	// }
	// fmt.Println("2. Добавление баннеров завершено")

	// // 3. Попытка получить баннер от имени пользователя с несовпадающим тегом+фичей
	// resp, err = client.GetUserBanner(baseURL, userToken, 2, 1)
	// if resp.StatusCode != http.StatusNotFound {
	// 	t.Fatalf("Ожидается статус код %d, получено %d", http.StatusNotFound, resp.StatusCode)
	// }
	// fmt.Println("3. Попытка получить баннер от пользователя с несовпадающим тегом успешно обработана")

	// // 4. Получение баннера от имени пользователя с совпадающим тегом
	// resp, err = client.GetUserBanner(baseURL, userToken, 1, 1)
	// if err != nil || resp.StatusCode != http.StatusOK {
	// 	t.Fatalf("Ожидается статус код %d, получено %d с ошибкой: %v", http.StatusOK, resp.StatusCode, err)
	// }
	// fmt.Println("4. Получение баннера от пользователя с совпадающим тегом успешно завершено")

	// // 5. Получение баннера от имени админа
	// resp, err = client.GetUserBanner(baseURL, userToken, 1, 1)
	// if err != nil || resp.StatusCode != http.StatusOK {
	// 	t.Fatalf("Ожидается статус код %d, получено %d с ошибкой: %v", http.StatusOK, resp.StatusCode, err)
	// }
	// fmt.Println("5. Получение баннера от пользователя с совпадающим тегом успешно завершено")

	// // 6. Попытка изменить баннер от пользователя без подходящего тега
	// resp, err = client.UpdateBanner(baseURL, userToken, 1, banner2)
	// if resp.StatusCode != http.StatusForbidden {
	// 	t.Fatalf("Ожидается статус код %d, получено %d", http.StatusForbidden, resp.StatusCode)
	// }
	// fmt.Println("6. Попытка изменить баннер от пользователя без подходящего тега успешно обработана")

	// // 7. Изменение баннера от админа
	// resp, err = client.UpdateBanner(baseURL, adminToken, 1, banner2)
	// if err != nil || resp.StatusCode != http.StatusOK {
	// 	t.Fatalf("Ожидается статус код %d, получено %d с ошибкой: %v", http.StatusOK, resp.StatusCode, err)
	// }
	// fmt.Println("7. Изменение баннера от пользователя с подходящим тегом успешно завершено")

	// 8. Проверка изменений баннера
	resp, err := client.GetUserBanner(baseURL, userToken, 2, 2)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидается статус код %d, получено %d с ошибкой: %v", http.StatusOK, resp.StatusCode, err)
	}
	///////////////////////////////////////////////////
	var updatedBanner model.Banner

	// Прочитать тело ответа
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Ошибка при чтении тела ответа: %v", err)
	}

	// Десериализовать тело ответа в структуру Banner
	if err := json.Unmarshal(body, &updatedBanner); err != nil {
		t.Fatalf("Ошибка при десериализации баннера из тела ответа: %v", err)
	}

	// Проверить, соответствует ли обновленное содержимое баннера содержимому banner2
	if !reflect.DeepEqual(updatedBanner.Content, banner2.Content) {
		t.Fatalf("Содержимое обновленного баннера не соответствует ожидаемому")
	}

	///////////////////////////////////////////////////
	fmt.Println("8. Проверка изменений баннера успешно завершена")

	resp, err = client.GetBanners(baseURL, userToken, 1, 1)
	if err == nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Ожидается ошибка доступа при выполнении GetBanners от лица пользователя, но ошибка не возникла")
	}

	// 9. Проверка изменений баннера
	resp, err = client.GetBanners(baseURL, adminToken, 1, 1)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидается статус код %d, получено %d с ошибкой: %v", http.StatusOK, resp.StatusCode, err)
	}
	///////////////////////////////////////////////////
	var updatedBanners []model.Banner

	// Прочитать тело ответа
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Ошибка при чтении тела ответа: %v", err)
	}

	// Десериализовать тело ответа в структуру Banner
	if err := json.Unmarshal(body, &updatedBanners); err != nil {
		t.Fatalf("Ошибка при десериализации баннера из тела ответа: %v", err)
	}

	// Проверить, соответствует ли обновленное содержимое баннера содержимому banner2
	if !reflect.DeepEqual(updatedBanners[0].Content, banner2.Content) {
		t.Fatalf("Содержимое обновленного баннера не соответствует ожидаемому")
	}

	///////////////////////////////////////////////////
	fmt.Println("9. Проверка изменений баннера успешно завершена")

	// 10. Удаление баннера от пользователя с несовпадающим тегом
	resp, err = client.DeleteBanner(baseURL, userToken, 1)
	if err == nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Ожидается статус код %d и ошибка, получено %d без ошибки", http.StatusForbidden, resp.StatusCode)
	}
	fmt.Println("10. Удаление баннера от пользователя с несовпадающим тегом успешно обработано")

	// 11. Удаление оставшихся двух баннеров от админа
	resp, err = client.DeleteBanner(baseURL, adminToken, 1)
	if err != nil || resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Ожидается статус код %d, получено %d с ошибкой: %v", http.StatusNoContent, resp.StatusCode, err)
	}
	resp, err = client.DeleteBanner(baseURL, adminToken, 2)
	if err != nil || resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Ожидается статус код %d, получено %d с ошибкой: %v", http.StatusNoContent, resp.StatusCode, err)
	}
	resp, err = client.DeleteBanner(baseURL, adminToken, 3)
	if err != nil || resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Ожидается статус код %d, получено %d с ошибкой: %v", http.StatusNoContent, resp.StatusCode, err)
	}
	fmt.Println("11. Удаление оставшихся двух баннеров от админа успешно завершено")

	// 12. Проверка отсутствия оставшихся баннеров
	resp, err = client.GetBanners(baseURL, adminToken, 1, 1)
	if err == nil || resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Ожидается статус код %d и ошибка, получено %d без ошибки", http.StatusNotFound, resp.StatusCode)
	}
	fmt.Println("12. Проверка отсутствия оставшихся баннеров успешно завершена")
}
