package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/15Andrew43/backend-trainee-assignment-2024/model"
)

func AddBanner(url string, token string, banner model.BannerRequest) (*http.Response, error) {
	body, err := json.Marshal(banner)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сериализации баннера: %v", err)
	}

	req, err := http.NewRequest("POST", url+"/banner", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("token", token)

	return http.DefaultClient.Do(req)
}

func UpdateBanner(url string, token string, id int, banner model.BannerRequest) (*http.Response, error) {
	body, err := json.Marshal(banner)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сериализации баннера: %v", err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/banner/%d", url, id), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("token", token)

	return http.DefaultClient.Do(req)
}

func GetBanners(url string, token string, tagID int, featureID int) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/banner?tag_id=%d&feature_id=%d", url, tagID, featureID), nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("token", token)

	return http.DefaultClient.Do(req)
}

func GetUserBanner(url string, token string, tagID int, featureID int, useLastRevision bool) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user_banner?tag_id=%d&feature_id=%d&use_last_revision=%v", url, tagID, featureID, useLastRevision), nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("token", token)

	return http.DefaultClient.Do(req)
}

func DeleteBanner(url string, token string, id int) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/banner/%d", url, id), nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("token", token)

	return http.DefaultClient.Do(req)
}
