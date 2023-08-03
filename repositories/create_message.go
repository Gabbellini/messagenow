package repositories

import "messagenow/domain/entities"

type CreateMessageRepository interface {
	Execute(roomID, senderID int64, message entities.Message) error
}
