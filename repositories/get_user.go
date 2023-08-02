package repositories

import (
	"context"
	"messagenow/domain/entities"
)

type GetUserRepository interface {
	Execute(ctx context.Context, userID int64) (*entities.User, error)
}
