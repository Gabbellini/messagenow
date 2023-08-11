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

func (c createMessageRepository) Execute(roomID, senderID int64, message entities.Message) error {
	//language=sql
	query := `
	INSERT INTO message (id_room, id_user, text)
	VALUES (?, ?, ?)`

	_, err := c.db.Exec(query, roomID, senderID, message.Text)
	if err != nil {
		log.Println("[createMessageRepository] Error Exec", err)
		return exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}

	return nil
}
