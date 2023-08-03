package usecases

import (
	"context"
	"messagenow/domain/entities"
)

type GetPreviousMessagesUseCase interface {
	Execute(ctx context.Context, userID, roomID int64) ([]entities.Message, error)
}
