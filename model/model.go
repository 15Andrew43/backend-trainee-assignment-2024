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

type RequestBodyBanner struct {
	TagIds    []int  `json:"tag_ids"`
	FeatureId int    `json:"feature_id"`
	Content   string `json:"content"`
	IsActive  bool   `json:"is_active"`
}
