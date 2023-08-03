package usecases

import (
	"messagenow/domain/entities"
)

type CreateTextMessageUseCase interface {
	Execute(Message entities.Message, senderID int64, addresseeID int64) error
}
