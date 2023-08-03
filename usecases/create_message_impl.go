package usecases

import (
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
	"messagenow/repositories"
)

type createMessageUseCaseImpl struct {
	createMessageRepository repositories.CreateMessageRepository
	getRoomByID             repositories.GetRoomRepository
}

func NewCreateMessageUseCase(
	createMessageRepository repositories.CreateMessageRepository,
	getRoomRepository repositories.GetRoomRepository,
) CreateMessageUseCase {
	return createMessageUseCaseImpl{
		createMessageRepository: createMessageRepository,
		getRoomByID:             getRoomRepository,
	}
}

func (c createMessageUseCaseImpl) Execute(senderID, roomID int64, message entities.Message) error {
	room, err := c.getRoomByID.Execute(roomID, senderID)
	if err != nil {
		log.Println("[createMessageUseCaseImpl] Error getRoomByID")
		return err
	}

	if room == nil {
		return exceptions.NewUnauthorizedError(exceptions.UnauthorizedMessage)
	}

	return c.createMessageRepository.Execute(roomID, senderID, message)
}
