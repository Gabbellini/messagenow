package repositories

import (
	"context"
	"messagenow/domain/entities"
)

type CreateRoomRepository interface {
	Execute(ctx context.Context, roomID int64, addresseeID int64) (*entities.Room, error)
}
