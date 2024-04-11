package database

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/15Andrew43/backend-trainee-assignment-2024/config"
	"github.com/15Andrew43/backend-trainee-assignment-2024/model"
	"go.mongodb.org/mongo-driver/bson"
)

func GetPostgresBanner(tagID, featureID int, banner *model.Banner) error {
	return PgConn.QueryRow(context.Background(), `
				SELECT b.id, b.data_id, b.is_active
				FROM banners b
				INNER JOIN banner_tags bt ON b.id = bt.banner_id
				WHERE b.feature_id = $1 AND bt.tag_id = $2
			`, featureID, tagID).Scan(&banner.ID, &banner.DataID, &banner.IsActive)
}

func GetMongoBannerData(bannerData *model.BannerData, banner *model.Banner) error {

	collection := MongoCli.Database(config.Cfg.MongoDB).Collection(config.Cfg.MongoCollection)

	////       TODO: strnig -> int        //////////////////////////////////////////////////////////////////////////
	dataID, err := strconv.Atoi(banner.DataID)
	if err != nil {
		log.Printf("ошибка преобразования строки в число: %v", err)
		return errors.New("can not convert str to int")
	}
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	filter := bson.M{"id": dataID}
	return collection.FindOne(context.Background(), filter).Decode(&bannerData)
}
