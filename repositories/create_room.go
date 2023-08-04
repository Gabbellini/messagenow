package repositories

import (
	"context"
	"messagenow/domain/entities"
)

type CreateRoomRepository interface {
	Execute(ctx context.Context) (*entities.Room, error)
}
