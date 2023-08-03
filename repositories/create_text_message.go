package repositories

import (
	"messagenow/domain/entities"
)

type CreateTextMessageRepository interface {
	Execute(Message entities.Message, senderID int64, addresseeID int64) error
}
