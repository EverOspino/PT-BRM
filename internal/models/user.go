package models

import (
	"errors"
	"regexp"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// ErrInvalidEmailFormat es retornado cuando el formato de un email es inválido.
var ErrInvalidEmailFormat = errors.New("el email no es válido")

// Validaciones de negocio
func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("el nombre es requerido")
	}
	if u.Email == "" {
		return errors.New("el email es requerido")
	}
	if !IsValidEmail(u.Email) {
		return ErrInvalidEmailFormat
	}
	if u.Age < 0 || u.Age > 150 {
		return errors.New("la edad debe estar entre 0 y 150")
	}
	return nil
}

func IsValidEmail(email string) bool {
	// Expresión regular simple para validar el formato del email. ejemplo: mail@mail.com
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}
