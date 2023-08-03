package repositories

import (
	"database/sql"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
)

type createMessageRepository struct {
	db *sql.DB
}

func NewCreateMessageRepository(db *sql.DB) CreateMessageRepository {
	return createMessageRepository{db: db}
}

func (c createMessageRepository) Execute(roomID, senderID, addresseeID int64, message entities.Message) error {
	//language=sql
	query := `
	INSERT INTO message (id_room, id_sender, id_addressee, text)
	VALUES (?, ?, ?, ?)`

	_, err := c.db.Exec(query, addresseeID, roomID, senderID, message.Text)
	if err != nil {
		log.Println("[createMessageRepository] Error Exec", err)
		return exceptions.NewInternalServerError(exceptions.InternalErrorMessage)
	}

	return nil
}
