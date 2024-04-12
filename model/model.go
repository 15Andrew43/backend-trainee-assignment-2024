package model

type PostgresBanner struct {
	ID       int    `json:"id"`
	DataID   string `json:"data_id"`
	IsActive bool   `json:"is_active"`
}

type BannerData map[string]interface{}

type Banner struct {
	TagIds    []int                  `json:"tag_ids"`
	FeatureId int                    `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
}
