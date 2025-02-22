package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	repo, err := NewRepository(os.Getenv("DATABASE_URI"))
	if err != nil {
		log.Fatal(err)
	}

	if err := repo.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewServer(repo)

	fmt.Println("We are flying...")
	server.Run()
}
