package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type IRepository interface {
	Save(url, code string) error
	Find(code string) (UrlObject, error)
	IncreaseAccessCount(code string) error
	Update(url, code string) error
	Delete(code string) error
}

type ShorteningServer struct {
	Repo IRepository
}

func (s *ShorteningServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/shorten", makeHttpHandlerFunc(s.shortHandler)).Methods(http.MethodPost)
	router.HandleFunc("/shorten/{code}", makeHttpHandlerFunc(s.getHandler)).Methods(http.MethodGet)
	router.HandleFunc("/shorten/{code}", makeHttpHandlerFunc(s.updateHandler)).Methods(http.MethodPut)
	router.HandleFunc("/shorten/{code}", makeHttpHandlerFunc(s.deleteHandler)).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":3000", router))
}

func NewServer(repository IRepository) *ShorteningServer {
	return &ShorteningServer{
		Repo: repository,
	}
}

func (s *ShorteningServer) shortHandler(w http.ResponseWriter, r *http.Request) error {

	var creationRequest CreationRequest

	if err := json.NewDecoder(r.Body).Decode(&creationRequest); err != nil {
		return err
	}

	_, err := url.ParseRequestURI(creationRequest.Url)
	if err != nil {
		return errInvalidURLFormat
	}

	code := generateShortCode()

	if err = s.Repo.Save(creationRequest.Url, code); err != nil {
		return err
	}
	return WriteJSON(w, 201, CreationRequestResponse{
		ShortCode: code,
	})
}

func (s *ShorteningServer) getHandler(w http.ResponseWriter, r *http.Request) error {

	code := mux.Vars(r)["code"]

	urlObj, err := s.Repo.Find(code)
	if err != nil {
		return err
	}

	go func() {
		if err := s.Repo.IncreaseAccessCount(code); err != nil {
			log.Printf("Failed to increase access count for code %s: %v", code, err)
		}
	}()
	return WriteJSON(w, http.StatusOK, urlObj)
}
func (s *ShorteningServer) deleteHandler(w http.ResponseWriter, r *http.Request) error {
	code := mux.Vars(r)["code"]

	err := s.Repo.Delete(code)
	return err
}

func (s *ShorteningServer) updateHandler(w http.ResponseWriter, r *http.Request) error {

	code := mux.Vars(r)["code"]

	var updationRequest UpdationRequest

	if err := json.NewDecoder(r.Body).Decode(&updationRequest); err != nil {
		return err
	}

	_, err := url.ParseRequestURI(updationRequest.Url)
	if err != nil {
		return errInvalidURLFormat
	}

	if err := s.Repo.Update(updationRequest.Url, code); err != nil {
		return err
	}

	return WriteJSON(w, 200, nil)
}

type APIFuncError struct {
	Error string `json:"error"`
}

func makeHttpHandlerFunc(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, GetStatusCode(err), APIFuncError{
				Error: err.Error(),
			})
		}
	}
}

func WriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	if v == nil {
		w.WriteHeader(statusCode)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

func GetStatusCode(err error) int {
	if code, exists := ErrStatusCodes[err]; exists {
		return code
	}
	return http.StatusInternalServerError
}

func generateShortCode() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 7)
	for i := range code {
		code[i] = letters[rand.IntN(len(letters))]
	}
	return string(code)
}
