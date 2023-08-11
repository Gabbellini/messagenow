package usecases

import (
	"context"
)

type CreateRoomUseCase interface {
	Execute(ctx context.Context) (int64, error)
}
