package repositories

import (
	"context"
	"messagenow/domain/entities"
)

type CreateRoomRepository interface {
	Execute(ctx context.Context, room entities.Room) (int64, error)
}
