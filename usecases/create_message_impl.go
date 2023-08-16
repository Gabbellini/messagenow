package usecases

import (
	"messagenow/domain/entities"
	"messagenow/repositories"
)

type createMessageUseCaseImpl struct {
	createMessageRepository repositories.CreateMessageRepository
}

func NewCreateMessageUseCase(
	createMessageRepository repositories.CreateMessageRepository,
) CreateMessageUseCase {
	return createMessageUseCaseImpl{
		createMessageRepository: createMessageRepository,
	}
}

func (c createMessageUseCaseImpl) Execute(senderID, roomID int64, message entities.Message) error {
	return c.createMessageRepository.Execute(roomID, senderID, message)
}
