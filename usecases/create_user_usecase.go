package usecases

import (
	"context"
	"messagenow/domain/entities"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, user entities.User) (int64, error)
}
