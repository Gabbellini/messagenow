package usecases

import "context"

type JoinRoomUseCase interface {
	Execute(ctx context.Context, roomID, userID int64) error
}
