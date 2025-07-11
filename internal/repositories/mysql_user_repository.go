package repositories

import (
	"database/sql"
	"fmt"
	"pt-brm/internal/database"
	"pt-brm/internal/models"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetAll() ([]*models.User, error)
	GetByID(id int) (*models.User, error)
	Update(id int, user *models.User) (*models.User, error)
	Delete(id int) error
	GetByEmail(email string) (*models.User, error)
}

type MySQLUserRepository struct {
	db *database.DB
}

func NewMySQLUserRepository(db *database.DB) UserRepository {
	return &MySQLUserRepository{
		db: db,
	}
}

func (r *MySQLUserRepository) Create(user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (name, email, age, created_at, updated_at) 
		VALUES (?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.Exec(query, user.Name, user.Email, user.Age)
	if err != nil {
		// Verificar si es error de email duplicado
		if isDuplicateKeyError(err) {
			return nil, fmt.Errorf("el email ya existe")
		}
		return nil, fmt.Errorf("no se pudo crear el usuario: %w", err)
	}

	// Obtener el ID generado
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener el id del registro insertado: %w", err)
	}

	// Retornar el usuario creado
	return r.GetByID(int(id))
}

func (r *MySQLUserRepository) GetAll() ([]*models.User, error) {
	query := `
		SELECT id, name, email, age, created_at, updated_at 
		FROM users 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("no se pudo consultar los usuarios: %w", err)
	}
	// Asegurarse de cerrar las filas al final
	defer rows.Close()

	var users []*models.User
	// Iterar sobre las filas obtenidas
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Age,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("no se pudo escanear el usuario: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar filas: %w", err)
	}

	return users, nil
}

func (r *MySQLUserRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, name, email, age, created_at, updated_at 
		FROM users 
		WHERE id = ?
	`

	user := &models.User{}
	// Ejecutar la consulta y escanear el resultado
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("no se pudo obtener al usuario: %w", err)
	}

	return user, nil
}

func (r *MySQLUserRepository) Update(id int, user *models.User) (*models.User, error) {
	query := `
		UPDATE users 
		SET name = ?, email = ?, age = ?, updated_at = NOW() 
		WHERE id = ?
	`

	result, err := r.db.Exec(query, user.Name, user.Email, user.Age, id)
	if err != nil {
		// Verificar si es error de email duplicado
		if isDuplicateKeyError(err) {
			return nil, fmt.Errorf("el email ya existe")
		}
		return nil, fmt.Errorf("no se pudo actualizar el usuario: %w", err)
	}

	// Verificar si se actualizó alguna fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("no se pudieron obtener las filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	// Retornar el usuario actualizado
	return r.GetByID(id)
}

func (r *MySQLUserRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el usuario: %w", err)
	}

	// Verificar si se eliminó alguna fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no se pudieron obtener las filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuario no encontrado")
	}

	return nil
}

func (r *MySQLUserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, name, email, age, created_at, updated_at 
		FROM users 
		WHERE email = ?
	`

	user := &models.User{}
	// Ejecutar la consulta y escanear el resultado
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("no se pudo obtener el usuario: %w", err)
	}

	return user, nil
}

// Helper function para detectar errores de clave duplicada en MySQL
func isDuplicateKeyError(err error) bool {
	// MySQL error code 1062 = Duplicate entry
	return err != nil &&
		(err.Error() == "Error 1062: Duplicate entry" ||
			fmt.Sprintf("%v", err)[:4] == "1062")
}
