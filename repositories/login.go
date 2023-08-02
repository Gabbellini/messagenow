package repositories

import (
	"context"
	"messagenow/domain/entities"
)

type LoginRepository interface {
	Execute(ctx context.Context, credential entities.Credential) (*entities.User, error)
}
