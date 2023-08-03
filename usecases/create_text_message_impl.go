package usecases

import (
	"messagenow/domain/entities"
	"messagenow/repositories"
)

type createTextMessageUseCaseImpl struct {
	createTextMessageRepository repositories.CreateTextMessageRepository
}

func NewCreateTextMessageUseCase(createTextMessageRepository repositories.CreateTextMessageRepository) CreateTextMessageUseCase {
	return createTextMessageUseCaseImpl{createTextMessageRepository: createTextMessageRepository}
}

func (c createTextMessageUseCaseImpl) Execute(messageText entities.MessageText, senderID int64, addresseeID int64) error {
	return c.createTextMessageRepository.Execute(messageText, senderID, addresseeID)
}
