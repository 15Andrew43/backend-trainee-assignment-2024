package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Banner struct {
	ID       int    `json:"id"`
	DataID   string `json:"data_id"`
	IsActive bool   `json:"is_active"`
}

type BannerData struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

var db *pgx.Conn
var dbMongo *mongo.Client

var (
	pgHost     string
	pgPort     int
	pgUser     string
	pgPassword string
	pgDB       string

	mongoHost       string
	mongoPort       int
	mongoDB         string
	mongoCollection string
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	pgHost := os.Getenv("POSTGRES_CONTAINER_NAME")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDB := os.Getenv("POSTGRES_DB")

	mongoHost = os.Getenv("MONGO_CONTAINER_NAME")
	mongoPort := os.Getenv("MONGO_PORT")
	mongoDB = os.Getenv("MONGO_DB")
	mongoCollection = os.Getenv("MONGO_COLLECTION")

	pgConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pgUser, pgPassword, pgHost, pgPort, pgDB)
	conn, err := pgx.Connect(context.Background(), pgConnString)
	if err != nil {
		log.Fatal("ошибка подключения к Postgres: %v", err)
	}
	defer conn.Close(context.Background())
	db = conn
	if err := db.PgConn().Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("ошибка подключения к MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())
	dbMongo = client
	if err := dbMongo.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", getUserBanner)

	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func getPostgresBanner(tagID, featureID int, banner *Banner) error {
	err := db.QueryRow(context.Background(), `
	SELECT b.id, b.data_id, b.is_active
	FROM banners b
	INNER JOIN banner_tags bt ON b.id = bt.banner_id
	WHERE b.feature_id = $1 AND bt.tag_id = $2
`, featureID, tagID).Scan(&banner.ID, &banner.DataID, &banner.IsActive)
	return err
}

func getMongoBannerData(bannerData *BannerData, banner *Banner) error {
	collection := dbMongo.Database(mongoDB).Collection(mongoCollection)

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

func getUserBanner(w http.ResponseWriter, r *http.Request) {
	// postgres
	tagID, _ := strconv.Atoi(r.URL.Query().Get("tag_id"))
	featureID, _ := strconv.Atoi(r.URL.Query().Get("feature_id"))

	var banner Banner
	err := getPostgresBanner(tagID, featureID, &banner)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("no rows with tag_id = %d and feature_id = %d", tagID, featureID)
			http.Error(w, "Баннер для не найден", http.StatusNotFound)
			return
		}
		log.Printf("ошибка при выполнении запроса к Postgres: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// mongo
	var bannerData BannerData
	err = getMongoBannerData(&bannerData, &banner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("не найдено документов с data_id = %v", banner.DataID)
			http.Error(w, "Баннер не найден", http.StatusNotFound)
			return
		}
		log.Printf("ошибка при выполнении запроса к MongoDB: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bannerData)
}

// middleware : token chech -> validate input data
