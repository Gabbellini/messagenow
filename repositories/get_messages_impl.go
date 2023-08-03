package repositories

import (
	"context"
	"database/sql"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
)

type getMessagesRepositoryImpl struct {
	db *sql.DB
}

func NewGetMessagesRepository(db *sql.DB) GetMessagesRepository {
	return getMessagesRepositoryImpl{
		db: db,
	}
}

func (g getMessagesRepositoryImpl) Execute(ctx context.Context, userID, roomID int64) ([]entities.Message, error) {
	//language=sql
	query := `
	SELECT u_sender.id, 
	       u_sender.name, 
	       u_sender.image,
	       u_addressee.id, 
	       u_addressee.name, 
	       u_addressee.image,
	       m.text
	FROM message m
	         INNER JOIN room r on m.id_room = r.id
	         INNER JOIN user_room ur on r.id = ur.id_room AND ur.id_user = ?
			 INNER JOIN user u_sender on m.id_sender = u_sender.id
			 INNER JOIN user u_addressee on m.id_addressee = u_addressee.id
 	WHERE m.id_room = ?`

	rows, err := g.db.QueryContext(ctx, query, userID, roomID)
	if err != nil {
		log.Println("[getMessagesRepositoryImpl] Error QueryContext", err)
		return nil, exceptions.NewInternalServerError(exceptions.InternalErrorMessage)
	}
	defer rows.Close()

	var messages []entities.Message
	for rows.Next() {
		var message entities.Message
		err = rows.Scan(
			&message.Sender.ID,
			&message.Sender.Name,
			&message.Sender.Image,
			&message.Addressee.ID,
			&message.Addressee.Name,
			&message.Addressee.Image,
			&message.Text,
		)
		if err != nil {
			log.Println("[getMessagesRepositoryImpl] Error ExecContext", err)
			return nil, exceptions.NewInternalServerError(exceptions.InternalErrorMessage)
		}

		messages = append(messages, message)
	}

	return messages, nil
}
