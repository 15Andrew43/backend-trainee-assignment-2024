package model

type Banner struct {
	ID       int    `json:"id"`
	DataID   string `json:"data_id"`
	IsActive bool   `json:"is_active"`
}

type BannerData struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}
