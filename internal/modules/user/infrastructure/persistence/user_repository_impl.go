package persistence

import (
	"go-hexagonal-template/internal/modules/user/domain/model"
	"go-hexagonal-template/internal/modules/user/domain/port"

	"gorm.io/gorm"
)

// DBInterface define la interfaz para las operaciones de base de datos
type DBInterface interface {
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
}

// UserRepositoryImpl implementa la interfaz UserRepository
type UserRepositoryImpl struct {
	db DBInterface
}

// NewUserRepositoryImpl crea una nueva instancia de UserRepositoryImpl
func NewUserRepositoryImpl(db DBInterface) port.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// Create implementa el método Create de la interfaz UserRepository
func (r *UserRepositoryImpl) Create(user *model.User) (*model.User, error) {
	result := r.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// GetByID implementa el método GetByID de la interfaz UserRepository
func (r *UserRepositoryImpl) GetByID(id string) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail implementa el método GetByEmail de la interfaz UserRepository
func (r *UserRepositoryImpl) GetByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
