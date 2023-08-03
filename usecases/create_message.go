package usecases

import "messagenow/domain/entities"

type CreateMessageUseCase interface {
	Execute(userID, roomID int64, message entities.Message) error
}
