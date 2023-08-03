package usecases

import (
	"context"
	"messagenow/domain/entities"
	"messagenow/repositories"
)

type getMessagesUseCase struct {
	getMessagesRepository repositories.GetMessagesRepository
}

func NewGetMessagesUseCase(getMessagesRepository repositories.GetMessagesRepository) GetMessagesUseCase {
	return getMessagesUseCase{
		getMessagesRepository: getMessagesRepository,
	}
}

func (g getMessagesUseCase) Execute(ctx context.Context, userID, roomID int64) ([]entities.Message, error) {
	return g.getMessagesRepository.Execute(ctx, userID, roomID)
}
