package repositories

import (
	"context"
	"messagenow/domain/entities"
)

type GetRoomRepository interface {
	Execute(ctx context.Context, roomID int64, userID int64) (*entities.Room, error)
}
