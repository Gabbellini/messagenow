package repositories

import (
	"context"
	"messagenow/domain/entities"
)

type GetPreviousMessagesRepository interface {
	Execute(ctx context.Context, userID, roomID int64) ([]entities.Message, error)
}
