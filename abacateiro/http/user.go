package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterUserRoutes(router chi.Router) {
	router.Get("/users", s.handleGetUsers)
	router.Get("/users/{id}", s.handleGetUser)
	router.Post("/users", s.handleCreateUser)
	router.Put("/users/{id}", s.handleUpdateUser)
	router.Delete("/users/{id}", s.handleDeleteUser)
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	// var user core.User	
	log.Println("GET /users/{id}")
}

func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	// var filter core.UserFilterDTO
	s.userService.GetUsers()
	log.Println("GET /users")
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	// var user core.User
	log.Println("POST /users")
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// var user core.UserUpdateDTO
	log.Println("PUT /users/{id}")
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	// var user core.User
	log.Println("DELETE /users/{id}")
}

