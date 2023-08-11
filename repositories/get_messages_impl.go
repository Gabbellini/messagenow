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
	SELECT u.id, 
		   u.name, 
		   u.image,
		   m.text,
		   m.created_at
	FROM message m
			 INNER JOIN room r on m.id_room = r.id
			 INNER JOIN user_room ur on r.id = ur.id_room AND ur.id_user = ?
			 INNER JOIN user u on u.id = m.id_user
	WHERE m.id_room = ?
	ORDER BY created_at`

	rows, err := g.db.QueryContext(ctx, query, userID, roomID)
	if err != nil {
		log.Println("[getMessagesRepositoryImpl] Error QueryContext", err)
		return nil, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}
	defer rows.Close()

	var messages []entities.Message
	for rows.Next() {
		var message entities.Message
		err = rows.Scan(
			&message.Sender.ID,
			&message.Sender.Name,
			&message.Sender.ImageURL,
			&message.Text,
			&message.CreatedAt,
		)
		if err != nil {
			log.Println("[getMessagesRepositoryImpl] Error ExecContext", err)
			return nil, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
		}

		messages = append(messages, message)
	}

	return messages, nil
}
