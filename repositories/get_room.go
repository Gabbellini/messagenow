package repositories

import (
	"messagenow/domain/entities"
)

type GetRoomRepository interface {
	Execute(roomID int64, senderID int64) (*entities.Room, error)
}
