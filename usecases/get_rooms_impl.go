package usecases

import (
	"context"
	"messagenow/domain/entities"
	"messagenow/repositories"
)

type getRoomsUseCaseImpl struct {
	getRoomRepository repositories.GetRoomsRepository
}

func NewGetRoomsUseCase(getRoomRepository repositories.GetRoomsRepository) GetRoomsUseCase {
	return getRoomsUseCaseImpl{
		getRoomRepository: getRoomRepository,
	}
}

func (g getRoomsUseCaseImpl) Execute(ctx context.Context, userID int64) ([]entities.Room, error) {
	return g.getRoomRepository.Execute(ctx, userID)
}
