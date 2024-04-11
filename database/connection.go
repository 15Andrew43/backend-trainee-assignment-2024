package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/15Andrew43/backend-trainee-assignment-2024/config"
)

var (
	PgConn   *pgx.Conn
	MongoCli *mongo.Client
)

func ConnectToPostgres(cfg *config.Config) error {
	pgConnString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.PGUser, cfg.PGPassword, cfg.PGHost, cfg.PGPort, cfg.PGDB)
	conn, err := pgx.Connect(context.Background(), pgConnString)
	if err != nil {
		return fmt.Errorf("ошибка подключения к Postgres: %v", err)
	}
	PgConn = conn
	return PgConn.Ping(context.Background())
}

func ConnectToMongoDB(cfg *config.Config) error {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", cfg.MongoHost, cfg.MongoPort))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("ошибка подключения к MongoDB: %v", err)
	}
	MongoCli = client
	return MongoCli.Ping(context.Background(), nil)
}
