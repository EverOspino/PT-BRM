package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
func LoadConfig() (*Config, error) {
	// Cargar variables desde archivo .env si existe
	if err := loadEnvFile(".env"); err != nil {
		// Si no existe .env, solo log warning pero continúa
		fmt.Printf("Warning: .env file not found, using environment variables or defaults\n")
	}

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

// loadEnvFile carga variables de entorno desde un archivo .env
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Saltar líneas vacías y comentarios
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Dividir en clave=valor
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Solo establecer si no existe ya en el entorno
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// Obtiene el valor de una variable de entorno o devuelve un valor por defecto si no está definida.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
