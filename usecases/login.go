package usecases

import (
	"context"
	"messagenow/domain/entities"
)

type LoginUseCase interface {
	Execute(ctx context.Context, credential entities.Credentials) (*entities.User, error)
}
