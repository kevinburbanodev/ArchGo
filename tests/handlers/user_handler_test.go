package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-hexagonal-template/internal/handlers"
	"go-hexagonal-template/internal/modules/user/application"
	"go-hexagonal-template/internal/modules/user/domain/model"
	"go-hexagonal-template/tests/mocks"

	"go-hexagonal-template/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *mocks.UserRepositoryMock) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockRepo := mocks.NewUserRepositoryMock()
	userHandler := handlers.NewUserHandler(mockRepo)
	router.POST("/users", userHandler.CreateUser)
	router.POST("/login", userHandler.Login)
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware()) // Usar el middleware real
	api.GET("/users/:id", userHandler.GetUser)
	return router, mockRepo
}

func TestUserHandler_GetUser(t *testing.T) {
	// Arrange
	router, mockRepo := setupTestRouter()

	// Crear un usuario de prueba
	now := time.Now()
	user := &model.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "Test",
		CreatedAt: now,
		UpdatedAt: now,
	}
	createdUser, err := mockRepo.Create(user)
	assert.NoError(t, err, "Error al crear el usuario para la prueba")
	assert.NotNil(t, createdUser, "El usuario creado no debería ser nil")
	assert.NotEmpty(t, createdUser.ID, "El ID del usuario no debería estar vacío")

	// Login para obtener el token
	loginInput := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	loginJson, _ := json.Marshal(loginInput)
	loginReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJson))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)

	var loginResponse struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(loginW.Body.Bytes(), &loginResponse)
	assert.NoError(t, err, "Error al deserializar la respuesta del login")
	assert.NotEmpty(t, loginResponse.Token, "El token no debería estar vacío")

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users/"+fmt.Sprintf("%d", createdUser.ID), bytes.NewBuffer(nil))
	req.Header.Set("Authorization", "Bearer "+loginResponse.Token)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code, "El código de estado debería ser 200")
	assert.NotEmpty(t, w.Body.String(), "El cuerpo de la respuesta no debería estar vacío")

	var response model.User
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Error al deserializar la respuesta")
	assert.Equal(t, createdUser.ID, response.ID, "El ID no coincide")
	assert.Equal(t, createdUser.Email, response.Email, "El email no coincide")
	assert.Equal(t, createdUser.Name, response.Name, "El nombre no coincide")
	assert.WithinDuration(t, createdUser.CreatedAt, response.CreatedAt, time.Second, "La fecha de creación no coincide")
	assert.WithinDuration(t, createdUser.UpdatedAt, response.UpdatedAt, time.Second, "La fecha de actualización no coincide")
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	// Arrange
	router, _ := setupTestRouter()

	// Login para obtener el token
	loginInput := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	loginJson, _ := json.Marshal(loginInput)
	loginReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJson))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)

	var loginResponse struct {
		Token string `json:"token"`
	}
	err := json.Unmarshal(loginW.Body.Bytes(), &loginResponse)
	assert.NoError(t, err, "Error al deserializar la respuesta del login")
	assert.NotEmpty(t, loginResponse.Token, "El token no debería estar vacío")

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users/9999", bytes.NewBuffer(nil))
	req.Header.Set("Authorization", "Bearer "+loginResponse.Token)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code, "El código de estado debería ser 404")
	assert.NotEmpty(t, w.Body.String(), "El cuerpo de la respuesta no debería estar vacío")

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Error al deserializar la respuesta")
	assert.Contains(t, response, "error", "La respuesta debería contener un campo 'error'")
	assert.Equal(t, "Usuario no encontrado", response["error"], "El mensaje de error no coincide")
}

func TestUserHandler_CreateUser(t *testing.T) {
	// Arrange
	router, _ := setupTestRouter()
	input := application.CreateUserInput{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "password123",
	}
	jsonInput, _ := json.Marshal(input)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code, "El código de estado debería ser 201")
	assert.NotEmpty(t, w.Body.String(), "El cuerpo de la respuesta no debería estar vacío")

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Error al deserializar la respuesta")
	assert.NotEmpty(t, response.ID, "El ID no debería estar vacío")
	assert.Equal(t, input.Email, response.Email, "El email no coincide")
	assert.Equal(t, input.Name, response.Name, "El nombre no coincide")
	assert.Empty(t, response.Password, "La contraseña no debería estar en la respuesta")
	assert.NotEmpty(t, response.CreatedAt, "La fecha de creación no debería estar vacía")
	assert.NotEmpty(t, response.UpdatedAt, "La fecha de actualización no debería estar vacía")
}

func TestUserHandler_CreateUser_InvalidInput(t *testing.T) {
	// Arrange
	router, _ := setupTestRouter()
	input := application.CreateUserInput{
		Email:    "invalid-email", // Email inválido
		Name:     "",              // Nombre vacío
		Password: "123",           // Contraseña muy corta
	}
	jsonInput, _ := json.Marshal(input)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code, "El código de estado debería ser 400")
	assert.NotEmpty(t, w.Body.String(), "El cuerpo de la respuesta no debería estar vacío")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Error al deserializar la respuesta")
	assert.Contains(t, response, "error", "La respuesta debería contener un campo 'error'")
	assert.Equal(t, "Datos de usuario inválidos", response["error"], "El mensaje de error no coincide")
}

func TestUserHandler_Login(t *testing.T) {
	// Arrange
	router, _ := setupTestRouter()
	input := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonInput, _ := json.Marshal(input)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code, "El código de estado debería ser 200")
	assert.NotEmpty(t, w.Body.String(), "El cuerpo de la respuesta no debería estar vacío")

	var response struct {
		Token string     `json:"token"`
		User  model.User `json:"user"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Error al deserializar la respuesta")
	assert.NotEmpty(t, response.Token, "El token no debería estar vacío")
	assert.NotEmpty(t, response.User.ID, "El ID del usuario no debería estar vacío")
	assert.Equal(t, input["email"], response.User.Email, "El email no coincide")
	assert.Empty(t, response.User.Password, "La contraseña no debería estar en la respuesta")
}

func TestUserHandler_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	router, _ := setupTestRouter()
	input := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "wrongpassword",
	}
	jsonInput, _ := json.Marshal(input)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code, "El código de estado debería ser 401")
	assert.NotEmpty(t, w.Body.String(), "El cuerpo de la respuesta no debería estar vacío")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Error al deserializar la respuesta")
	assert.Contains(t, response, "error", "La respuesta debería contener un campo 'error'")
	assert.Equal(t, "Credenciales inválidas", response["error"], "El mensaje de error no coincide")
}

func TestUserHandler_GetUser_WithAuth(t *testing.T) {
	// Arrange
	router, mockRepo := setupTestRouter()

	// Crear un usuario de prueba
	now := time.Now()
	user := &model.User{
		ID:        2,
		Email:     "test@example.com",
		Name:      "Test",
		CreatedAt: now,
		UpdatedAt: now,
	}
	createdUser, err := mockRepo.Create(user)
	assert.NoError(t, err, "Error al crear el usuario para la prueba")
	assert.NotNil(t, createdUser, "El usuario creado no debería ser nil")
	assert.NotEmpty(t, createdUser.ID, "El ID del usuario no debería estar vacío")

	// Login para obtener el token
	loginInput := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	loginJson, _ := json.Marshal(loginInput)
	loginReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJson))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)

	var loginResponse struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(loginW.Body.Bytes(), &loginResponse)
	assert.NoError(t, err, "Error al deserializar la respuesta del login")
	assert.NotEmpty(t, loginResponse.Token, "El token no debería estar vacío")

	// Act - Obtener usuario con el token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users/"+fmt.Sprintf("%d", createdUser.ID), bytes.NewBuffer(nil))
	req.Header.Set("Authorization", "Bearer "+loginResponse.Token)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code, "El código de estado debería ser 200")
	assert.NotEmpty(t, w.Body.String(), "El cuerpo de la respuesta no debería estar vacío")

	var response model.User
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Error al deserializar la respuesta")
	assert.Equal(t, createdUser.ID, response.ID, "El ID no coincide")
	assert.Equal(t, createdUser.Email, response.Email, "El email no coincide")
	assert.Equal(t, createdUser.Name, response.Name, "El nombre no coincide")
	assert.WithinDuration(t, createdUser.CreatedAt, response.CreatedAt, time.Second, "La fecha de creación no coincide")
	assert.WithinDuration(t, createdUser.UpdatedAt, response.UpdatedAt, time.Second, "La fecha de actualización no coincide")
}

func TestUserHandler_GetUser_WithoutAuth(t *testing.T) {
	// Arrange
	router, _ := setupTestRouter()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users/9999", bytes.NewBuffer(nil))
	// No se añade el token de autorización
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code, "El código de estado debería ser 401")
	assert.NotEmpty(t, w.Body.String(), "El cuerpo de la respuesta no debería estar vacío")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Error al deserializar la respuesta")
	assert.Contains(t, response, "error", "La respuesta debería contener un campo 'error'")
	assert.Equal(t, "No se proporcionó token de autenticación", response["error"], "El mensaje de error no coincide")
}
