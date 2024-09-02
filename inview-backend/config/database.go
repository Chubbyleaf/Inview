package config

import (
	"context"
	"fmt"
	"insense-local/database"
	"log"
	"time"

	"github.com/chenyahui/gin-cache/persist"
	"github.com/go-redis/redis/v8"
)

func NewMongoDatabase(env *Env) database.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass
	dbName := env.DBName

	mongodbUri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	if dbUser == "" || dbPass == "" {
		mongodbUri = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}

	// log.Printf("MongoDB URI : %s", mongodbUri)

	client, err := database.NewClient(mongodbUri)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func CloseMongoDBConnection(client database.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}

func NewRedisConnection(env *Env) *redis.Client {
	redisHost := env.RedisHost
	redisPort := env.RedisPort

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:   0,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Redis ping response: %s", pong)

	return client
}

type RedisCache struct {
	Store            *persist.RedisStore
	DefaultCacheTime time.Duration
}

func SetupRedisCache(redisHost string, redisPort string, password string) *RedisCache {
	return &RedisCache{
		Store: persist.NewRedisStore(redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr: fmt.Sprintf(
				"%s:%s",
				redisHost,
				redisPort,
			),
		})),
		DefaultCacheTime: 10 * time.Second,
	}
}
