package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/15Andrew43/backend-trainee-assignment-2024/config"
)

var (
	PgPool      *pgxpool.Pool
	MongoCli    *mongo.Client
	RedisClient *redis.Client
)

func ConnectToPostgres(cfg *config.Config) error {

	connURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?pool_max_conns=%d",
		cfg.PGUser, cfg.PGPassword, cfg.PGHost, cfg.PGPort, cfg.PGDB, 30)

	poolConfig, err := pgxpool.ParseConfig(connURL)
	if err != nil {
		return fmt.Errorf("ошибка парсинга строки подключения: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return fmt.Errorf("ошибка подключения к Postgres: %v", err)
	}

	PgPool = pool
	return PgPool.Ping(context.Background())
}

func ConnectToMongoDB(cfg *config.Config) error {
	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%d", cfg.MongoHost, cfg.MongoPort)).
		SetMaxPoolSize(30)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("ошибка подключения к MongoDB: %v", err)
	}

	MongoCli = client
	return MongoCli.Ping(context.Background(), nil)
}

func ConnectToRedis(cfg *config.Config) error {
	redisOptions := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		PoolSize: 30,
	}

	client := redis.NewClient(redisOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("ошибка подключения к Redis: %v", err)
	}

	RedisClient = client
	return nil
}
