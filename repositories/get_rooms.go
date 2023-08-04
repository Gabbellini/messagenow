package repositories

import (
	"context"
	"messagenow/domain/entities"
)

type GetRoomsRepository interface {
	Execute(ctx context.Context, userID int64) ([]entities.Room, error)
}
