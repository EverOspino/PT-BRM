package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// Carga la configuración desde las variables de entorno y devuelve una instancia de Config.
func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			Database: getEnv("DB_NAME", "database"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}, nil
}

// Genera una cadena de conexión para la base de datos.
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}

// Obtiene el valor de una variable de entorno o devuelve un valor por defecto si no está definida.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
