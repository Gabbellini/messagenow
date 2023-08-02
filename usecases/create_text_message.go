package usecases

import (
	"context"
	"messagenow/domain/entities"
)

type CreateTextMessageUseCase interface {
	Execute(context context.Context, messageText entities.MessageText, senderID int64, addresseeID int64) error
}
