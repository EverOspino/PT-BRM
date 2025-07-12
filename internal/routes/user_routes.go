package routes

import (
	"pt-brm/internal/handlers"

	"github.com/gorilla/mux"
)

// SetupUserRoutes configura las rutas específicas de usuarios
func SetupUserRoutes(router *mux.Router, userHandler *handlers.UserHandler) {
	users := router.PathPrefix("/users").Subrouter()

	users.HandleFunc("", userHandler.CreateUser).Methods("POST")
	users.HandleFunc("", userHandler.GetAllUsers).Methods("GET")
	users.HandleFunc("/{id}", userHandler.GetUserByID).Methods("GET")
	users.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
	users.HandleFunc("/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Rutas adicionales específicas de usuarios
	users.HandleFunc("/email/{email}", userHandler.GetByEmail).Methods("GET")
}
