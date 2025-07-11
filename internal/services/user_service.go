package services

import (
	"pt-brm/internal/models"
	"pt-brm/internal/repositories"
)

type UserService interface {
	CreateUser(req *models.CreateUserRequest) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetUserByID(id int) (*models.User, error)
	UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id int) error
	GetUserByEmail(email string) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// crear un nuevo usuario a partir de la solicitud
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	// Validar el usuario antes de crear
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Crear el usuario en el repositorio
	return s.userRepo.Create(user)
}

func (s *userService) GetAllUsers() ([]*models.User, error) {
	// Obtener todos los usuarios del repositorio
	return s.userRepo.GetAll()
}

func (s *userService) GetUserByID(id int) (*models.User, error) {
	// Obtener un usuario por ID del repositorio
	return s.userRepo.GetByID(id)
}

func (s *userService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error) {
	// Obtener el usuario existente por ID
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar los campos del usuario con los datos de la solicitud
	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age

	// Validar el usuario actualizado
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Actualizar el usuario en el repositorio
	return s.userRepo.Update(id, user)
}

func (s *userService) DeleteUser(id int) error {
	// Verificar si el usuario existe antes de eliminar
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Eliminar el usuario por ID del repositorio
	return s.userRepo.Delete(id)
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	// Validar el formato del email antes de buscar
	if !models.IsValidEmail(email) {
		return nil, models.ErrInvalidEmailFormat
	}

	// Obtener el usuario por email del repositorio
	return s.userRepo.GetByEmail(email)
}
