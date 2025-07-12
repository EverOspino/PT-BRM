package routes

import (
	"net/http"
	"pt-brm/internal/database"

	"github.com/gorilla/mux"
)

// SetupHealthRoutes configura las rutas de health check
func SetupHealthRoutes(router *mux.Router, db *database.DB) {
	router.HandleFunc("/health", healthCheck(db)).Methods("GET")
	router.HandleFunc("/ping", pingHandler()).Methods("GET")
}

func healthCheck(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "no se pudo conectar con la base de datos", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "database": "connected"}`))
	}
}

func pingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "pong"}`))
	}
}
