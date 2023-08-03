package repositories

import (
	"messagenow/domain/entities"
)

type GetRoomRepository interface {
	Execute(roomID int64, addreseeID int64) (*entities.Room, error)
}
