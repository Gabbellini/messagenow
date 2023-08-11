package usecases

import (
	"context"
	"messagenow/domain/entities"
)

type GetUserByIDUseCase interface {
	Execute(ctx context.Context, userID int64) (*entities.User, error)
}
