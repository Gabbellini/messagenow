package usecases

import "messagenow/domain/entities"

type CreateMessageUseCase interface {
	Execute(roomID, userID, addresseeID int64, message entities.Message) error
}
