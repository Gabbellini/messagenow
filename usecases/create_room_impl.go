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

func (c createRoomUseCasesImpl) Execute(ctx context.Context, room entities.Room) (int64, error) {
	err := c.processRoom(&room)
	if err != nil {
		log.Println("[createRoomUseCasesImpl] Error processRoom", err)
		return 0, err
	}

	return c.createRoomRepository.Execute(ctx, room)
}

func (c createRoomUseCasesImpl) processRoom(room *entities.Room) error {
	if room.Name = strings.TrimSpace(room.Name); room.Name == "" {
		return exceptions.NewBadRequestError("Nome da sala não pode ser vazio")
	}

	return nil
}
