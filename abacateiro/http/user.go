package http

import (
	"application"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// func dump(data interface{}) {
//     jsonData, err := json.MarshalIndent(data, "", "  ")
// 	if err != nil {
// 		fmt.Println("Erro ao serializar dados:", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println(string(jsonData))
// }

// func dd(data interface{}) {
//     dump(data)
// 	os.Exit(0)
// }

type errorResponse struct {
    Error string `json:"error"`
}

func (s *Server) RegisterUserRoutes(router chi.Router) {
	router.Get("/users", s.handleGetUsers)
	router.Get("/users/{id}", s.handleGetUser)
	router.Post("/users", s.handleCreateUser)
	router.Put("/users/{id}", s.handleUpdateUser)
	router.Delete("/users/{id}", s.handleDeleteUser)
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Invalid user ID"})
		return
	}
	
    user, err := s.userService.GetUser(id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
    users, err := s.userService.GetUsers()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }

    userDTOs := []application.UserResponse{}
    if len(users) > 0 {
        userDTOs = application.ToUserResponses(users)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(userDTOs)
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {

    var user application.User

    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: "Invalid request payload"})
        return
    }

    createdUser, err := s.userService.CreateUser(user)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdUser)
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {

    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: "Invalid user ID"})
        return
    }

    var user application.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: "Invalid request payload"})
        return
    }

    user.ID = id

    updatedUser, err := s.userService.UpdateUser(user)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedUser)
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: "Invalid user ID"})
        return
    }

    if err := s.userService.DeleteUser(id); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
