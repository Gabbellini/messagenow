package usecases

import (
	"context"
	"messagenow/domain/entities"
)

type AddUserRoomUseCase interface {
	Execute(ctx context.Context, user entities.User, roomID, userID int64) error
}
