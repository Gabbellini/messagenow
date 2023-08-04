package usecases

import (
	"context"
	"messagenow/domain/entities"
)

type CreateRoomUseCase interface {
	Execute(ctx context.Context) (*entities.Room, error)
}
