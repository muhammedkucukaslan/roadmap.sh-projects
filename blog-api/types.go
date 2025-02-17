package main

import (
	"errors"
	"net/http"
	"time"
)

type Blog struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// request body structs
type CreateBlogRequest struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

type UpdateBlogRequest struct {
	Title    *string   `json:"title,omitempty"` // Pointer to allow optional updates
	Content  *string   `json:"content,omitempty"`
	Category *string   `json:"category,omitempty"`
	Tags     *[]string `json:"tags,omitempty"`
}

// Errors
var (
	errServerError        = errors.New("interval server error")
	errInvalidRequestBody = errors.New("invalid request body")
	errInvalidID          = errors.New("invalid id")
)

type ErrorStatusCodes map[error]int

var ErrStatusCodes = ErrorStatusCodes{
	errInvalidRequestBody: http.StatusBadRequest,
	errInvalidRequestBody: http.StatusBadRequest,
	errInvalidID:          http.StatusBadRequest,
}
