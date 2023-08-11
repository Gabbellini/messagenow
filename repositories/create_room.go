package repositories

import (
	"context"
)

type CreateRoomRepository interface {
	Execute(ctx context.Context) (*int64, error)
}
