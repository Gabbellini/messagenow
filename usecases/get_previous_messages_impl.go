package usecases

import (
	"context"
	"messagenow/domain/entities"
	"messagenow/repositories"
)

type getPreviousMessagesRepositoryUseCase struct {
	getPreviousMessagesRepository repositories.GetPreviousMessagesRepository
}

func NewGetPreviousMessagesRepository(getPreviousMessagesRepository repositories.GetPreviousMessagesRepository) GetPreviousMessagesUseCase {
	return getPreviousMessagesRepositoryUseCase{
		getPreviousMessagesRepository: getPreviousMessagesRepository,
	}
}

func (g getPreviousMessagesRepositoryUseCase) Execute(ctx context.Context, userID, roomID int64) ([]entities.Message, error) {
	return g.getPreviousMessagesRepository.Execute(ctx, userID, roomID)
}
