package repositories

import (
	"context"
	"messagenow/domain/entities"
)

type CreateTextMessageRepository interface {
	Execute(context context.Context, messageText entities.MessageText, senderID int64, addresseeID int64) error
}
