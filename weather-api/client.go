package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type WeatherData struct {
	ResolvedAdress string `json:"resolvedAddress"`
	Timezone       string `json:"timezone"`
	Days           []struct {
		Datetime    string  `json:"datetime"`
		Sunrise     string  `json:"sunrise"`
		Sunset      string  `json:"sunset"`
		Tempmax     float64 `json:"tempmax"`
		Tempmin     float64 `json:"tempmin"`
		Temp        float64 `json:"temp"`
		Description string  `json:"description"`
	} `json:"days"`
}

type ErrorStruct struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var (
	serverError = ErrorStruct{
		Message: "SERVER ERROR",
		Code:    500,
	}
)

type Client struct {
	Ctx         context.Context
	redisClient *redis.Client
}

func NewClient(ctx context.Context, redisClient *redis.Client) *Client {
	return &Client{
		Ctx:         ctx,
		redisClient: redisClient,
	}
}

func (c *Client) weatherAPIHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	city := parts[len(parts)-1]

	if weatherData, err := c.cacheFromRedisDB(city); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(weatherData)
		return
	}

	weatherData, errStruct := c.fetchFromAPI(city)
	if errStruct != (ErrorStruct{}) {
		if errStruct.Code == http.StatusBadRequest {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errStruct)
		return
	}

	if err := c.setToRedisDB(city, weatherData); err != nil {
		logError(err.Error())
	}
	json.NewEncoder(w).Encode(weatherData)

}

func (c *Client) cacheFromRedisDB(key string) (WeatherData, error) {
	weatherData, err := c.redisClient.Get(c.Ctx, key).Result()
	if err == redis.Nil {
		return WeatherData{}, fmt.Errorf("key does not exist")
	} else if err != nil {
		logError(err.Error())
		return WeatherData{}, fmt.Errorf("error retrieving from Redis: %v", err)
	}

	var data WeatherData
	err = json.Unmarshal([]byte(weatherData), &data)
	if err != nil {
		logError(err.Error())
		return WeatherData{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return data, nil
}

func (c *Client) fetchFromAPI(city string) (WeatherData, ErrorStruct) {
	url := formatUrl(city)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)
	if err != nil {
		return WeatherData{}, serverError
	}
	defer res.Body.Close()
	var weatherData WeatherData

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {

			return WeatherData{}, serverError
		}
		return WeatherData{}, ErrorStruct{
			Message: string(bodyBytes),
			Code:    400,
		}
	}

	if err = json.NewDecoder(res.Body).Decode(&weatherData); err != nil {
		logError(err.Error())
		return WeatherData{}, serverError
	}

	return weatherData, ErrorStruct{}
}

func (c *Client) setToRedisDB(city string, weatherData WeatherData) error {
	jsonData, err := json.Marshal(weatherData)
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}
	return c.redisClient.Set(c.Ctx, city, string(jsonData), 24*time.Hour).Err()
}
func formatUrl(city string) string {
	api_key := os.Getenv("WEATHER_API_KEY")
	return fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?unitGroup=metric&include=stats%%2Cfcst&key=%s&contentType=json", city, api_key)
}

func logError(errMsg string) {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	report := fmt.Sprintf("%s : %s", time.Now().Format(time.RFC3339), errMsg)
	log.Println(report)
}
