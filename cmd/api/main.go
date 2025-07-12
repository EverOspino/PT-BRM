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
	"pt-brm/internal/routes"
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

	// Crear el router
	router := routes.NewRouter(db)
	// Configurar las rutas
	httpHandler := router.SetupRoutes()

	// Iniciar el servidor HTTP
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      httpHandler,
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
