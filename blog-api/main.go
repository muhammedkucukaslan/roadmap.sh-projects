package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	repo, err := NewRepository(os.Getenv("DATABASE_URI"))

	if err != nil {
		log.Fatal(err)
	}

	if err := repo.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":3000", repo)

	fmt.Println("We are flying")
	server.Run()
}
