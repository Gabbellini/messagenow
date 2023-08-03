package usecases

import (
	"messagenow/domain/entities"
)

type CreateTextMessageUseCase interface {
	Execute(messageText entities.MessageText, senderID int64, addresseeID int64) error
}
