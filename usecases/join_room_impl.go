package usecases

import (
	"context"
	"messagenow/repositories"
)

type joinRoomUseCaseImpl struct {
	joinRoomRepository repositories.JoinRoomRepository
}

func NewJoinRoomUseCase(joinRoomRepository repositories.JoinRoomRepository) JoinRoomUseCase {
	return joinRoomUseCaseImpl{
		joinRoomRepository: joinRoomRepository,
	}
}

func (j joinRoomUseCaseImpl) Execute(ctx context.Context, roomID, userID int64) error {
	return j.joinRoomRepository.Execute(ctx, roomID, userID)
}
