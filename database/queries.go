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
	"go.mongodb.org/mongo-driver/bson"
)

func GetPostgresBanner(tagID, featureID int, banner *model.PostgresBanner) error {
	return PgPool.QueryRow(context.Background(), `
				SELECT b.id, b.data_id, b.is_active
				FROM banners b
				INNER JOIN banner_tags bt ON b.id = bt.banner_id
				WHERE b.feature_id = $1 AND bt.tag_id = $2
			`, featureID, tagID).Scan(&banner.ID, &banner.DataID, &banner.IsActive)
}

func GetPostgresAllBanners(tagID, featureID, limit, offset int) ([]model.PostgresBanner, error) {
	rows, err := PgPool.Query(context.Background(), `
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

	var banners []model.PostgresBanner
	for rows.Next() {
		var banner model.PostgresBanner
		if err := rows.Scan(&banner.ID, &banner.DataID, &banner.IsActive); err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}
	return banners, nil
}

func CreatePostgresBanner(nextId int, ch chan<- error, requestBody *model.Banner) {

	// check that banners with such feature + tag do not exist
	for _, tag := range requestBody.TagIds {
		var banner model.PostgresBanner
		err := GetPostgresBanner(tag, requestBody.FeatureId, &banner)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				continue
			}
		}
		ch <- my_errors.BannerExist{Feature_id: requestBody.FeatureId, Tag_id: tag}
		return
	}

	var insertedID int
	err := PgPool.QueryRow(context.Background(), `
					INSERT INTO banners (feature_id, data_id, is_active)
					VALUES ($1, $2, $3)
					RETURNING id;
				`, requestBody.FeatureId, strconv.Itoa(nextId), requestBody.IsActive).Scan(&insertedID)
	if err != nil {
		ch <- err
		return
	}
	log.Printf("Вставлена новая строка %v в таблицу banners", insertedID)

	for _, tag := range requestBody.TagIds {
		_, err = PgPool.Exec(context.Background(), `
					INSERT INTO banner_tags (banner_id, tag_id)
					VALUES ($1, $2);
				`, insertedID, tag)
		if err != nil {
			ch <- err
			return
		}
	}
	log.Printf("Вставлены новые строки в таблицу banner_tags")

	ch <- nil
}

func UpgradePostgresBanner(id int, requestBody *model.Banner) (int, error) {

	// check that banners with such feature + tag do not exist
	for _, tag := range requestBody.TagIds {
		var banner model.PostgresBanner
		err := GetPostgresBanner(tag, requestBody.FeatureId, &banner)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				continue
			}
		}
		if banner.ID == id {
			continue
		}
		return 0, &my_errors.BannerExist{Feature_id: requestBody.FeatureId, Tag_id: tag}
	}

	var dataIdStr string
	err := PgPool.QueryRow(context.Background(), `
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

	_, err = PgPool.Exec(context.Background(), `
					UPDATE banners
					SET feature_id = $2, is_active = $3, updated_at = NOW()
					WHERE id = $1;
				`, id, requestBody.FeatureId, requestBody.IsActive)
	if err != nil {
		return 0, err
	}
	log.Printf("Произведено обновление содержимого баннера в таблице banners")

	_, err = PgPool.Exec(context.Background(), `
					DELETE FROM banner_tags
					WHERE banner_id = $1;
				`, id)
	if err != nil {
		return 0, err
	}
	log.Printf("При обновлении удалены строки из таблицы banner_tags")

	for _, tag := range requestBody.TagIds {
		_, err = PgPool.Exec(context.Background(), `
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

func DeletePostgresBanner(id int) (int, error) {
	var dataIdStr string
	err := PgPool.QueryRow(context.Background(), `
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

	_, err = PgPool.Exec(context.Background(), `
					DELETE FROM banner_tags
					WHERE banner_id = $1
				`, id)
	if err != nil {
		return 0, err
	}
	log.Printf("Удалены строки из таблицы banner_tags для баннера %v", id)

	_, err = PgPool.Exec(context.Background(), `
					DELETE FROM banners
					WHERE id = $1
				`, id)
	if err != nil {
		return 0, err
	}
	log.Printf("Удалена строка %v из таблицы banners", id)

	return dataId, nil
}

func GetMongoBannerData(bannerData *model.MongoBannerData, banner *model.PostgresBanner) error {

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

func CreateMongoBanner(nextId int, ch chan<- error, content map[string]interface{}) {

	collection := MongoCli.Database(config.Cfg.MongoDB).Collection(config.Cfg.MongoCollection)

	_, err := collection.InsertOne(context.Background(), map[string]interface{}{
		"id":      nextId,
		"content": content,
	})
	if err != nil {
		ch <- err
		return
	}

	ch <- nil
}

func UpgradeMongoBanner(dataId int, content map[string]interface{}) error {
	collection := MongoCli.Database(config.Cfg.MongoDB).Collection(config.Cfg.MongoCollection)

	filter := bson.M{"id": dataId}
	update := bson.M{"$set": bson.M{"content": content}}

	_, err := collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteMongoBanner(dataId int) error {
	collection := MongoCli.Database(config.Cfg.MongoDB).Collection(config.Cfg.MongoCollection)

	_, err := collection.DeleteOne(
		context.Background(),
		bson.M{"id": dataId},
	)
	if err != nil {
		return err
	}

	return nil
}
