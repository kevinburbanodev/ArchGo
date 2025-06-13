package handlers

import (
	"net/http"
	"strconv"

	"go-hexagonal-template/internal/modules/user/application"
	"go-hexagonal-template/internal/modules/user/domain/port"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	getUserUseCase    *application.GetUserUseCase
	createUserUseCase *application.CreateUserUseCase
	loginUserUseCase  *application.LoginUserUseCase
}

func NewUserHandler(userRepository port.UserRepository) *UserHandler {
	return &UserHandler{
		getUserUseCase:    application.NewGetUserUseCase(userRepository),
		createUserUseCase: application.NewCreateUserUseCase(userRepository),
		loginUserUseCase:  application.NewLoginUserUseCase(userRepository),
	}
}

// GetUser godoc
// @Summary Obtener usuario por ID
// @Description Obtiene los detalles de un usuario por su ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Security Bearer
// @Success 200 {object} model.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	idUint, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	user, err := h.getUserUseCase.Execute(uint(idUint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuario no encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary Crear un nuevo usuario
// @Description Crea un nuevo usuario en el sistema
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.User true "Datos del usuario"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var input application.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de usuario inválidos",
		})
		return
	}

	createdUser, err := h.createUserUseCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al crear el usuario",
		})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// Login godoc
// @Summary Autenticar usuario
// @Description Autentica un usuario y devuelve un token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body map[string]string true "Credenciales de usuario"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Credenciales inválidas",
		})
		return
	}

	result, err := h.loginUserUseCase.Execute(application.LoginUserInput{
		Email:    credentials.Email,
		Password: credentials.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Credenciales inválidas",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
