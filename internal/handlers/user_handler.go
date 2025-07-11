package handlers

import (
	"encoding/json"
	"net/http"
	"pt-brm/internal/models"
	"pt-brm/internal/services"
	"pt-brm/pkg/response"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// POST /users - Crear usuario
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	// Decodificar el cuerpo de la solicitud en la estructura CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

// GET /users - Obtener todos los usuarios
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, users)
}

// GET /users/{id} - Obtener usuario por ID
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario de los parámetros de la ruta
	vars := mux.Vars(r)
	// Convertir el ID de string a int
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	// Obtener el email de los parámetros de la ruta
	vars := mux.Vars(r)
	email := vars["email"]

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// PUT /users/{id} - Actualizar usuario
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario de los parámetros de la ruta
	vars := mux.Vars(r)
	// Convertir el ID de string a int
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// Decodificar el cuerpo de la solicitud en la estructura UpdateUserRequest
	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	user, err := h.userService.UpdateUser(id, &req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// DELETE /users/{id} - Eliminar usuario
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario de los parámetros de la ruta
	vars := mux.Vars(r)
	// Convertir el ID de string a int
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
