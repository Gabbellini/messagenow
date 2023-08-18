package usecases

import (
	"context"
	"messagenow/domain/entities"
)

type CreateRoomUseCase interface {
	Execute(ctx context.Context, user entities.User, room entities.Room) (*entities.Room, error)
}
