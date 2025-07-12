package routes

import (
	"net/http"
	"pt-brm/internal/database"
	"pt-brm/internal/handlers"
	"pt-brm/internal/repositories"
	"pt-brm/internal/services"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Router struct {
	db *database.DB
}

func NewRouter(db *database.DB) *Router {
	return &Router{db: db}
}

func (rt *Router) SetupRoutes() http.Handler {
	// Crear dependencias
	userRepo := repositories.NewMySQLUserRepository(rt.db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Router principal
	router := mux.NewRouter()

	// Rutas de salud
	SetupHealthRoutes(router, rt.db)

	// API v1
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Rutas por módulo
	SetupUserRoutes(apiV1, userHandler)

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	return c.Handler(router)
}
