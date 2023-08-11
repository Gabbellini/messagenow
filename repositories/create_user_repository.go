package repositories

import (
	"context"
	"messagenow/domain/entities"
)

type CreateUserRepository interface {
	Execute(ctx context.Context, user entities.User) (int64, error)
}
