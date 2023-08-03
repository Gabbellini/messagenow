package usecases

import (
	"context"
	"messagenow/domain/entities"
)

type CreateRoomUseCase interface {
	Execute(ctx context.Context, roomID int64, addresseeID int64) (*entities.Room, error)
}
