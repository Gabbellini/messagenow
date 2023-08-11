package repositories

import (
	"context"
	"messagenow/domain/entities"
)

type GetUserByEmailRepository interface {
	Execute(ctx context.Context, email string) (*entities.User, error)
}
