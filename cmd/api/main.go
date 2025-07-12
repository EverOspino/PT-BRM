package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pt-brm/internal/config"
	"pt-brm/internal/database"
	"pt-brm/internal/handlers"
	"pt-brm/internal/repositories"
	"pt-brm/internal/services"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Cargar la configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("no se pudo cargar la confuguracion de la base de datos: %v", err)
	}

	// Conectar a la base de datos
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Conexión a la base de datos exitosa")

	// Ejecutar migraciones
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Crear el repositorio de usuarios
	userRepo := repositories.NewMySQLUserRepository(db)

	// Crear el servicio de usuarios
	userService := services.NewUserService(userRepo)

	// Crear el manejador de usuarios
	userHandler := handlers.NewUserHandler(userService)

	// Configurar el router
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", healthCheck(db)).Methods("GET")

	// Definir las rutas para el manejador de usuarios
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	corsHandler := c.Handler(router)

	// Iniciar el servidor HTTP
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      corsHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Servidor corriendo en el puerto %s", cfg.Server.Port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Manejar señales de terminación
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("señal recibida %s, apagando el servidor...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("No se pudo apagar el servidor: %v", err)
	}

	log.Println("Servidor apagado correctamente")
}

// Health check endpoint
func healthCheck(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "Database connection failed", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "database": "connected"}`))
	}
}
