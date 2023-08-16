package repositories

import "context"

type SetRoomStatusRepository interface {
	Execute(ctx context.Context, roomID, status int64) error
}
