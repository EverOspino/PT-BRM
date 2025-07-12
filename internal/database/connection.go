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
		return nil, fmt.Errorf("error al comprobar Ping con la base de datos: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(80) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		age INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_email (email),
		INDEX idx_created_at (created_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("no se pudo crear la tabla users: %w", err)
	}

	return nil
}
