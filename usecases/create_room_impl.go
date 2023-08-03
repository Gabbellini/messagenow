package usecases

import (
	"context"
	"messagenow/domain/entities"
	"messagenow/repositories"
)

type createRoomUseCasesImpl struct {
	createRoomRepository repositories.CreateRoomRepository
}

func NewCreateRoomUseCase(createRoomRepository repositories.CreateRoomRepository) CreateRoomUseCase {
	return createRoomUseCasesImpl{
		createRoomRepository: createRoomRepository,
	}
}

func (c createRoomUseCasesImpl) Execute(ctx context.Context, userID int64, addresseeID int64) (*entities.Room, error) {
	return c.createRoomRepository.Execute(ctx, userID, addresseeID)
}
