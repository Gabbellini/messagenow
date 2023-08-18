package usecases

import (
	"context"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
	"messagenow/repositories"
	"strings"
)

type createRoomUseCasesImpl struct {
	createRoomRepository repositories.CreateRoomRepository
	getRoomRepository    repositories.GetRoomRepository
}

func NewCreateRoomUseCase(
	createRoomRepository repositories.CreateRoomRepository,
	getRoomRepository repositories.GetRoomRepository,
) CreateRoomUseCase {
	return createRoomUseCasesImpl{
		createRoomRepository: createRoomRepository,
		getRoomRepository:    getRoomRepository,
	}
}

func (c createRoomUseCasesImpl) Execute(ctx context.Context, user entities.User, room entities.Room) (*entities.Room, error) {
	err := c.processRoom(&room)
	if err != nil {
		log.Println("[createRoomUseCasesImpl] Error processRoom", err)
		return nil, err
	}

	roomID, err := c.createRoomRepository.Execute(ctx, room)
	if err != nil {
		log.Println("[createRoomUseCasesImpl] Error createRoomRepository", err)
		return nil, err
	}

	return c.getRoomRepository.Execute(ctx, user.ID, roomID)
}

func (c createRoomUseCasesImpl) processRoom(room *entities.Room) error {
	if room.Name = strings.TrimSpace(room.Name); room.Name == "" {
		return exceptions.NewBadRequestError("Nome da sala não pode ser vazio")
	}

	return nil
}
