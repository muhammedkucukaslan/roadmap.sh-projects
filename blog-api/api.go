package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type IRepoistory interface {
	Create(CreateBlogRequest) error
	Update(int, UpdateBlogRequest) error
	Delete(id int) error
	GetByID(id int) (Blog, error)
	GetBlogs() ([]Blog, error)
}

type APIServer struct {
	ListenAdrr string
	Repo       IRepoistory
}

func NewAPIServer(listenAdrr string, repo IRepoistory) *APIServer {
	return &APIServer{
		ListenAdrr: listenAdrr,
		Repo:       repo,
	}
}

func (s *APIServer) Run() {

	router := mux.NewRouter()

	router.HandleFunc("/blogs", makeHttpHandlerFunc(s.createBlogHandle)).Methods("POST")
	router.HandleFunc("/blogs", makeHttpHandlerFunc(s.getBlogsHandle)).Methods("GET")
	router.HandleFunc("/blogs/{id}", makeHttpHandlerFunc(s.getBlogByIDHandle)).Methods("GET")
	router.HandleFunc("/blogs/{id}", makeHttpHandlerFunc(s.updateBlogHandle)).Methods("PUT")
	router.HandleFunc("/blogs/{id}", makeHttpHandlerFunc(s.deleteBlogHandle)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(s.ListenAdrr, router))
}

func (s *APIServer) createBlogHandle(w http.ResponseWriter, r *http.Request) error {
	var createBlogRequest CreateBlogRequest

	if err := json.NewDecoder(r.Body).Decode(&createBlogRequest); err != nil {
		return errInvalidRequestBody
	}

	if err := s.Repo.Create(createBlogRequest); err != nil {
		return errServerError
	}

	w.WriteHeader(201)
	return nil
}

func (s *APIServer) deleteBlogHandle(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := s.Repo.Delete(id); err != nil {
		return err
	}
	w.WriteHeader(200)
	return nil
}

func (s *APIServer) getBlogByIDHandle(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	blog, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, 200, blog)
}

func (s *APIServer) getBlogsHandle(w http.ResponseWriter, r *http.Request) error {
	blogs, err := s.Repo.GetBlogs()
	if err != nil {
		return err
	}
	return WriteJSON(w, 200, blogs)
}

func (s *APIServer) updateBlogHandle(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var updateBlogRequest UpdateBlogRequest

	//validate request body

	if err := json.NewDecoder(r.Body).Decode(&updateBlogRequest); err != nil {
		return errInvalidRequestBody
	}

	if err := s.Repo.Update(id, updateBlogRequest); err != nil {
		return err
	}

	w.WriteHeader(200)

	return nil
}

type APIFuncError struct {
	Error string `json:"Error"`
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHttpHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, GetStatusCode(err), APIFuncError{
				Error: err.Error(),
			})
		}
	}
}

func WriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

func GetStatusCode(err error) int {
	if code, exists := ErrStatusCodes[err]; exists {
		return code
	}
	return http.StatusInternalServerError // Default fallback
}
