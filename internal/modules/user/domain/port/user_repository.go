package port

import "go-hexagonal-template/internal/modules/user/domain/model"

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	GetByID(id string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}
