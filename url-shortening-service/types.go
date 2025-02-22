package main

import (
	"errors"
	"net/http"
)

type UrlObject struct {
	ID          string `json:"id"`
	Url         string `json:"url"`
	ShortCode   string `json:"shortCode"`
	AccessCount int    `json:"accessCount"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type CreationRequest struct {
	Url string `json:"url"`
}

type CreationRequestResponse struct {
	ShortCode string `json:"shortCode"`
}

type UpdationRequest struct {
	Url string `json:"url"`
}

var (
	errInvalidURLFormat   = errors.New("invalid url format")
	errServerError        = errors.New("internal server error")
	errInvalidRequestBody = errors.New("invalid request body")
	errInvalidID          = errors.New("invalid id")
	errInvalidCode        = errors.New("invalid code")
	errUrlNotFound        = errors.New("url not found")
)

type ErrorStatusCodes map[error]int

var ErrStatusCodes = ErrorStatusCodes{
	errUrlNotFound:        http.StatusNotFound,
	errInvalidURLFormat:   http.StatusBadRequest,
	errInvalidRequestBody: http.StatusBadRequest,
	errInvalidRequestBody: http.StatusBadRequest,
	errInvalidID:          http.StatusBadRequest,
	errInvalidCode:        http.StatusBadRequest,
}
