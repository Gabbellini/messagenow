package repositories

import (
	"context"
	"database/sql"
	"log"
	"messagenow/exceptions"
)

type setRoomStatusRepository struct {
	db *sql.DB
}

func NewSetRoomStatusRepository(db *sql.DB) SetRoomStatusRepository {
	return setRoomStatusRepository{db: db}
}

func (s setRoomStatusRepository) Execute(ctx context.Context, roomID, status int64) error {
	//language=sql
	query := `
	UPDATE room 
	SET status = ?
	WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query, status, roomID)
	if err != nil {
		log.Println("[setRoomStatusRepository] Error ExecContext", err)
		return exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}

	return nil
}
