package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env file could not loaded")
	}
	redisClient := initRedis(os.Getenv("REDIS_SERVER_URL"))
	ctx := context.Background()

	client := NewClient(ctx, redisClient)

	http.HandleFunc("/weather/", client.weatherAPIHandler)

	fmt.Println("Server starting...")
	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}

func initRedis(serverUrl string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     serverUrl,
		Password: "", // No password set
		DB:       0,  // Use default DB
	})
	return client
}
