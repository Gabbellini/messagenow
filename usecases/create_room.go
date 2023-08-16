package usecases

import (
	"context"
	"messagenow/domain/entities"
)

type CreateRoomUseCase interface {
	Execute(ctx context.Context, room entities.Room) (int64, error)
}
