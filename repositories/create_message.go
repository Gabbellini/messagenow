package repositories

import "messagenow/domain/entities"

type CreateMessageRepository interface {
	Execute(roomID, senderID, addresseeID int64, message entities.Message) error
}
