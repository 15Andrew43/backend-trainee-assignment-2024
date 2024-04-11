package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost string
	ServerPort int

	PGHost     string
	PGPort     int
	PGUser     string
	PGPassword string
	PGDB       string

	MongoHost       string
	MongoPort       int
	MongoDB         string
	MongoCollection string

	RedisHost string
	RedisPort int
}

var Cfg Config

func LoadConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return err
	}

	pgPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return err
	}

	mongoPort, err := strconv.Atoi(os.Getenv("MONGO_PORT"))
	if err != nil {
		return err
	}

	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		return err
	}

	Cfg = Config{
		ServerHost: os.Getenv("SERVER_CONTAINER_NAME"),
		ServerPort: serverPort,

		PGHost:     os.Getenv("POSTGRES_CONTAINER_NAME"),
		PGPort:     pgPort,
		PGUser:     os.Getenv("POSTGRES_USER"),
		PGPassword: os.Getenv("POSTGRES_PASSWORD"),
		PGDB:       os.Getenv("POSTGRES_DB"),

		MongoHost:       os.Getenv("MONGO_CONTAINER_NAME"),
		MongoPort:       mongoPort,
		MongoDB:         os.Getenv("MONGO_DB"),
		MongoCollection: os.Getenv("MONGO_COLLECTION"),

		RedisHost: os.Getenv("REDIS_CONTAINER_NAME"),
		RedisPort: redisPort,
	}
	return nil
}
