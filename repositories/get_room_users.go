package repositories

import (
	"context"
	"messagenow/domain/entities"
)

type GetRoomUsersRepository interface {
	Execute(ctx context.Context, roomID int64) ([]entities.User, error)
}
