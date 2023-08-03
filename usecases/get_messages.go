package usecases

import (
	"context"
	"messagenow/domain/entities"
)

type GetMessagesUseCase interface {
	Execute(ctx context.Context, userID, roomID int64) ([]entities.Message, error)
}
