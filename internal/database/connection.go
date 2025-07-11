package database

import (
	"database/sql"
	"fmt"
	"pt-brm/internal/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

// Crea una nueva conexión a la base de datos utilizando la configuración proporcionada.
func NewConnection(cfg config.DatabaseConfig) (*DB, error) {
	// Crea la conneción de la base de datos
	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("error al abrir la conexión con la base de datos: %w", err)
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verificar conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error al conprobar Ping con la base de datos: %w", err)
	}

	return &DB{db}, nil
}
