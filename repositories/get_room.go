package repositories

import (
	"context"
	"messagenow/domain/entities"
)

type GetRoomRepository interface {
	Execute(ctx context.Context, userID int64, roomID int64) (*entities.Room, error)
}
