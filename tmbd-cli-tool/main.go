package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var (
	validTypes = []string{
		"playing", "popular", "top", "upcoming",
	}
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	args := os.Args

	if err := ValidateArgs(args); err != nil {
		HandleError(err)
	}

	filter := flag.String("type", "", "The type of movies you want to see")
	flag.Parse()
	key := os.Getenv("API_KEY")
	movies, err := fetchMovies(*filter, key)
	if err != nil {
		color.New(color.FgHiRed).Println(err)
		return
	}
	displayMovies(movies)

}

type Movie struct {
	Rating           float64 `json:"vote_average"`
	Overview         string  `json:"overview"`
	Adult            bool    `json:"adult"`
	Title            string  `json:"title"`
	OriginalLanguage string  `json:"original_language"`
	ReleaseDate      string  `json:"release_date"`
}

type Movies struct {
	Results []Movie `json:"results"`
}

type Error struct {
	Message string `json:"status_message"`
}

func fetchMovies(filter string, key string) (Movies, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s", getProperParam(filter))

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", key)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return Movies{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var err Error
		json.NewDecoder(res.Body).Decode(&err)
		return Movies{}, fmt.Errorf("%s", err.Message)
	}

	var movies Movies
	if err = json.NewDecoder(res.Body).Decode(&movies); err != nil {
		return Movies{}, err
	}
	return movies, nil

}
