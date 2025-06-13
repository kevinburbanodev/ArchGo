package persistence_test

import (
	"testing"
	"time"

	"go-hexagonal-template/internal/modules/user/domain/model"
	"go-hexagonal-template/internal/modules/user/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB es un mock de la base de datos para testing
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func setupTestDB() *MockDB {
	return new(MockDB)
}

func TestUserRepositoryImpl_Create(t *testing.T) {
	// Arrange
	mockDB := setupTestDB()
	repo := persistence.NewUserRepositoryImpl(mockDB)
	now := time.Now()
	user := &model.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "Test",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Configurar el mock para Create
	mockDB.On("Create", user).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*model.User)
		// Asegurarse de que el ID se mantiene
		assert.Equal(t, uint(1), arg.ID, "El ID debería mantenerse")
	}).Return(&gorm.DB{})

	// Act
	createdUser, err := repo.Create(user)

	// Assert
	assert.NoError(t, err, "Error al crear el usuario")
	assert.NotEmpty(t, createdUser.ID, "El ID del usuario no debería estar vacío")
	assert.Equal(t, uint(1), createdUser.ID, "El ID debería ser el mismo")
	assert.Equal(t, user.Email, createdUser.Email, "El email no coincide")
	assert.Equal(t, user.Name, createdUser.Name, "El nombre no coincide")
	assert.Equal(t, user.CreatedAt, createdUser.CreatedAt, "La fecha de creación no coincide")
	assert.Equal(t, user.UpdatedAt, createdUser.UpdatedAt, "La fecha de actualización no coincide")

	mockDB.AssertExpectations(t)
}

func TestUserRepositoryImpl_GetByID(t *testing.T) {
	// Arrange
	mockDB := setupTestDB()
	repo := persistence.NewUserRepositoryImpl(mockDB)
	now := time.Now()
	user := &model.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "Test",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Configurar el mock para First
	mockDB.On("First", mock.Anything, []interface{}{"id = ?", user.ID}).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*model.User)
		*arg = *user
	}).Return(&gorm.DB{})

	// Act
	foundUser, err := repo.GetByID(user.ID)

	// Assert
	assert.NoError(t, err, "Error al buscar el usuario por ID")
	assert.NotNil(t, foundUser, "El usuario encontrado no debería ser nil")
	assert.Equal(t, user.ID, foundUser.ID, "El ID no coincide")
	assert.Equal(t, user.Email, foundUser.Email, "El email no coincide")
	assert.Equal(t, user.Name, foundUser.Name, "El nombre no coincide")
	assert.Equal(t, user.CreatedAt, foundUser.CreatedAt, "La fecha de creación no coincide")
	assert.Equal(t, user.UpdatedAt, foundUser.UpdatedAt, "La fecha de actualización no coincide")

	mockDB.AssertExpectations(t)
}

func TestUserRepositoryImpl_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockDB := setupTestDB()
	repo := persistence.NewUserRepositoryImpl(mockDB)

	// Configurar el mock para First con error
	mockDB.On("First", mock.Anything, []interface{}{"id = ?", uint(9999)}).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})

	// Act
	user, err := repo.GetByID(9999)

	// Assert
	assert.Error(t, err, "Debería retornar un error cuando el usuario no existe")
	assert.Nil(t, user, "El usuario debería ser nil cuando no existe")

	mockDB.AssertExpectations(t)
}
