package usecases

import (
	"context"
	"messagenow/domain/entities"
)

type GetUserUseCase interface {
	Execute(ctx context.Context, userID int64) (*entities.User, error)
}
