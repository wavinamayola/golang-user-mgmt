package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	api "github.com/wavinamayola/user-management/internal/services"
)

func SetupRoutes(r *mux.Router, api *api.Service) {
	r.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the home page. This is just sample"))
	}).Methods("GET")
	r.HandleFunc("/login", api.UserService.Create).Methods("POST")

	s := r.PathPrefix("/users").Subrouter()
	r.HandleFunc("/users", api.UserService.Create).Methods("POST")
	s.HandleFunc("/", api.UserService.Create).Methods("POST")
	s.HandleFunc("/{id}", api.UserService.Get).Methods("GET")
	s.HandleFunc("/{id}", api.UserService.Update).Methods("PUT")
	s.HandleFunc("/{id}", api.UserService.Delete).Methods("DELETE")
}
