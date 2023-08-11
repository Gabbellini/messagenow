package usecases

import (
	"context"
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

func (c createRoomUseCasesImpl) Execute(ctx context.Context) (int64, error) {
	return c.createRoomRepository.Execute(ctx)
}
