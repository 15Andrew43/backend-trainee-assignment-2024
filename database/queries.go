package database

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/15Andrew43/backend-trainee-assignment-2024/config"
	"github.com/15Andrew43/backend-trainee-assignment-2024/model"
	my_errors "github.com/15Andrew43/backend-trainee-assignment-2024/my_errors"
	"github.com/15Andrew43/backend-trainee-assignment-2024/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
				SELECT DISTINCT b.id, b.data_id, b.is_active
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

	// check that banners with such feature + tag do not exist
	for _, tag := range requestBody.TagIds {
		var banner model.Banner
		err := GetPostgresBanner(tag, requestBody.FeatureId, &banner)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				continue
			}
		}
		return 0, &my_errors.BannerExist{Feature_id: requestBody.FeatureId, Tag_id: tag}
	}

	nextId := util.GenerateNextId()
	var insertedID int
	err := PgConn.QueryRow(context.Background(), `
					INSERT INTO banners (feature_id, data_id, is_active)
					VALUES ($1, $2, $3)
					RETURNING id;
				`, requestBody.FeatureId, strconv.Itoa(nextId), requestBody.IsActive).Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	log.Printf("Вставлена новая строка %v в таблицу banners", insertedID)

	for _, tag := range requestBody.TagIds {
		_, err = PgConn.Exec(context.Background(), `
					INSERT INTO banner_tags (banner_id, tag_id)
					VALUES ($1, $2);
				`, insertedID, tag)
		if err != nil {
			return 0, err
		}
	}
	log.Printf("Вставлены новые строки в таблицу banner_tags")

	return nextId, nil
}

func UpgradePostgresBanner(id int, requestBody *model.RequestBodyBanner) (int, error) {
	var dataIdStr string
	err := PgConn.QueryRow(context.Background(), `
					SELECT data_id
					FROM banners
					WHERE id = $1
				`, id).Scan(&dataIdStr)
	if err != nil {
		return 0, err
	}
	log.Printf("Получен data_id обновляемого баннера %v", id)

	dataId, err := strconv.Atoi(dataIdStr)
	if err != nil {
		return 0, err
	}

	_, err = PgConn.Exec(context.Background(), `
					UPDATE banners
					SET feature_id = $2, is_active = $3, updated_at = NOW()
					WHERE id = $1;
				`, id, requestBody.FeatureId, requestBody.IsActive)
	if err != nil {
		return 0, err
	}
	log.Printf("Произведено обновление содержимого баннера в таблице banners")

	_, err = PgConn.Exec(context.Background(), `
					DELETE FROM banner_tags
					WHERE banner_id = $1;
				`, id)
	if err != nil {
		return 0, err
	}
	log.Printf("При обновлении удалены строки из таблицы banner_tags")

	for _, tag := range requestBody.TagIds {
		_, err = PgConn.Exec(context.Background(), `
			INSERT INTO banner_tags (banner_id, tag_id)
			VALUES ($1, $2);
		`, id, tag)
		if err != nil {
			return 0, err
		}
	}
	log.Printf("Вставлены новые строки в таблицу banner_tags")

	return dataId, nil
}

func DeletePostgresBanner() {

	// var tmp int
	// err = PgConn.QueryRow(context.Background(), `
	// 				DELETE
	// 				FROM banners
	// 				WHERE id = $1;
	// 			`, id).Scan(&tmp)
	// if err != nil {
	// 	return 0, err
	// }
	// log.Printf("При обновлении удалена строка %v из таблицы banners", id)

	// err = PgConn.QueryRow(context.Background(), `
	// 				DELETE
	// 				FROM banner_tags
	// 				WHERE banner_id = $1;
	// 			`, id).Scan(&tmp)
	// if err != nil {
	// 	return 0, err
	// }
	// log.Printf("При обновлении удалены строки из таблицы banner_tags")

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

func CreateMongoBanner(nextId int, content map[string]interface{}) error {

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

func UpgradeMongoBanner(dataId int, content map[string]interface{}) error {
	collection := MongoCli.Database(config.Cfg.MongoDB).Collection(config.Cfg.MongoCollection)

	_, err := collection.ReplaceOne(
		context.Background(),
		bson.M{"id": dataId},
		content,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	return nil
}
