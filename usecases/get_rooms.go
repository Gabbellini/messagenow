package usecases

import (
	"context"
	"messagenow/domain/entities"
)

type GetRoomsUseCase interface {
	Execute(ctx context.Context, userID int64) ([]entities.Room, error)
}
