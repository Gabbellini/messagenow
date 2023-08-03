package repositories

import (
	"messagenow/domain/entities"
)

type CreateTextMessageRepository interface {
	Execute(messageText entities.MessageText, senderID int64, addresseeID int64) error
}
