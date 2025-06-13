package application

import (
	"go-hexagonal-template/internal/modules/user/domain/model"
	"go-hexagonal-template/internal/modules/user/domain/port"
)

type GetUserUseCase struct {
	userRepository port.UserRepository
}

func NewGetUserUseCase(userRepository port.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *GetUserUseCase) Execute(id uint) (*model.User, error) {
	return uc.userRepository.GetByID(id)
}
