package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

var (
	UserNotFound      = fmt.Errorf("User not found. Please check the username")
	ErrorFetchingData = fmt.Errorf("Error fetching data")
	JsonParsingError  = fmt.Errorf("Json parsing error")
)

type Client struct {
	username string
}

func NewClient(username string) *Client {
	return &Client{username}
}

func (c *Client) makeRequest() ([]activity, error) {

	resp, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/events", c.username))
	if err != nil {
		log.Fatalf("Error:%v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return nil, UserNotFound
	}

	if resp.StatusCode != http.StatusOK {
		var error fetcherror
		if err := json.NewDecoder(resp.Body).Decode(&error); err != nil {
			return nil, JsonParsingError
		}
		return nil, errors.New(error.Message)
	}

	var activities []activity
	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return nil, JsonParsingError
	}
	return activities, nil
}

type activity struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Created_at string `json:"created_at"`
	Payload    struct {
		Commits []struct {
			Message string `json:"message"`
		}
		Action   string `json:"action"`
		Ref_Type string `json:"ref_type"`
	} `json:"payload"`
}

type fetcherror struct {
	Message string `json:"message"`
}
