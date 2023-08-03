package repositories

import (
	"context"
	"database/sql"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
)

type getPreviousMessagesRepositoryImpl struct {
	db *sql.DB
}

func NewGetPreviousMessagesRepository(db *sql.DB) GetPreviousMessagesRepository {
	return getPreviousMessagesRepositoryImpl{
		db: db,
	}
}

func (g getPreviousMessagesRepositoryImpl) Execute(ctx context.Context, userID, roomID int64) ([]entities.Message, error) {
	//language=sql
	query := `
	SELECT text
	FROM message m
	         INNER JOIN room r on m.id_room = r.id
	         INNER JOIN user_room ur on r.id = ur.id_room AND ur.id_user = ?
	WHERE m.id_room = ?`

	rows, err := g.db.QueryContext(ctx, query, userID, roomID)
	if err != nil {
		log.Println("[getPreviousMessagesRepositoryImpl] Error QueryContext", err)
		return nil, exceptions.NewInternalServerError(exceptions.InternalErrorMessage)
	}
	defer rows.Close()

	var messages []entities.Message
	for rows.Next() {
		var message entities.Message
		err = rows.Scan(&message.Text)
		if err != nil {
			log.Println("[getPreviousMessagesRepositoryImpl] Error ExecContext", err)
			return nil, exceptions.NewInternalServerError(exceptions.InternalErrorMessage)
		}

		messages = append(messages, message)
	}

	return messages, nil
}
