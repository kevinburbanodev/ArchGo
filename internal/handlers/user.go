package handlers

import (
	"net/http"

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

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.getUserUseCase.Execute(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuario no encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		LastName string `json:"last_name" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de usuario inválidos",
		})
		return
	}

	user, err := h.createUserUseCase.Execute(application.CreateUserInput{
		Email:    input.Email,
		Name:     input.Name,
		LastName: input.LastName,
		Password: input.Password,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al crear el usuario",
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de login inválidos",
		})
		return
	}

	result, err := h.loginUserUseCase.Execute(application.LoginUserInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Credenciales inválidas",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
