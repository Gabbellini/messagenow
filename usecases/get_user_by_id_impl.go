package usecases

import (
	"context"
	"messagenow/domain/entities"
	"messagenow/repositories"
)

type getUserByIDUseCaseImpl struct {
	getUserByIDRepository repositories.GetUserByIDRepository
}

func NewGetUserByIDUseCase(getUserRepository repositories.GetUserByIDRepository) GetUserByIDUseCase {
	return getUserByIDUseCaseImpl{getUserByIDRepository: getUserRepository}
}

func (g getUserByIDUseCaseImpl) Execute(ctx context.Context, userID int64) (*entities.User, error) {
	return g.getUserByIDRepository.Execute(ctx, userID)
}
