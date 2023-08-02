package usecases

import (
	"context"
	"messagenow/domain/entities"
	"messagenow/repositories"
)

type getUserUseCaseImpl struct {
	getUserRepository repositories.GetUserRepository
}

func NewGetUserUseCase(getUserRepository repositories.GetUserRepository) GetUserUseCase {
	return getUserUseCaseImpl{getUserRepository: getUserRepository}
}

func (g getUserUseCaseImpl) Execute(ctx context.Context, userID int64) (*entities.User, error) {
	return g.getUserRepository.Execute(ctx, userID)
}
