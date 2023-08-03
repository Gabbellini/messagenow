package repositories

import "context"

type JoinRoomRepository interface {
	Execute(ctx context.Context, roomID, userID int64) error
}
