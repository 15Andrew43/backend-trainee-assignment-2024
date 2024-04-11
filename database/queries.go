package database

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/15Andrew43/backend-trainee-assignment-2024/config"
	"github.com/15Andrew43/backend-trainee-assignment-2024/model"
	"github.com/15Andrew43/backend-trainee-assignment-2024/util"
	"go.mongodb.org/mongo-driver/bson"
)

var UnsuccessfulInsert = errors.New("Unsuccessful Insert into Postgres")

func GetPostgresBanner(tagID, featureID int, banner *model.Banner) error {
	return PgConn.QueryRow(context.Background(), `
				SELECT b.id, b.data_id, b.is_active
				FROM banners b
				INNER JOIN banner_tags bt ON b.id = bt.banner_id
				WHERE b.feature_id = $1 AND bt.tag_id = $2
			`, featureID, tagID).Scan(&banner.ID, &banner.DataID, &banner.IsActive)
}

func GetPostgresAllBanners(tagID, featureID, limit, offset int) ([]model.Banner, error) {
	rows, err := PgConn.Query(context.Background(), `
				SELECT b.id, b.data_id, b.is_active
				FROM banners b
				INNER JOIN banner_tags bt ON b.id = bt.banner_id
				WHERE ($1 = -1 or b.feature_id = $1) AND ($2 = -1 OR bt.tag_id = $2)
				LIMIT $3
				OFFSET $4
			`, featureID, tagID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []model.Banner
	for rows.Next() {
		var banner model.Banner
		if err := rows.Scan(&banner.ID, &banner.DataID, &banner.IsActive); err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}
	return banners, nil
}

func CreatePostgresBanner(requestBody *model.RequestBodyBanner) (int, error) {
	nextId := util.GenerateNextId()
	var insertedID int

	log.Println("olololoolololololool")

	res, err := PgConn.Exec(context.Background(), `
					INSERT INTO banners (feature_id, data_id, is_active)
					VALUES ($1, $2, $3)
					RETURNING id;
				`, requestBody.FeatureId, strconv.Itoa(nextId), requestBody.IsActive)
	//.Scan(&insertedID)
	log.Printf("res = %+v\n\n\n", res)
	if err != nil {
		return 0, err
	}
	log.Println("qqwqwqwqwqwqwqqwqwqw")

	log.Println("kekekekekekekeke")

	for _, tag := range requestBody.TagIds {
		err = PgConn.QueryRow(context.Background(), `
					INSERT INTO banner_tags (banner_id, tag_id)
					VALUES ($1, $2);
				`, insertedID, tag).Scan(&insertedID)
		if err != nil {
			return 0, err
		}
	}

	return nextId, nil
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

func CreateMongoBanner(nextId int, content string) error {

	collection := MongoCli.Database(config.Cfg.MongoDB).Collection(config.Cfg.MongoCollection)

	_, err := collection.InsertOne(context.Background(), map[string]interface{}{
		"id":      nextId,
		"content": content,
	})
	if err != nil {
		return err
	}

	return nil
}
