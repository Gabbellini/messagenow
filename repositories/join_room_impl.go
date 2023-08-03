package repositories

import (
	"context"
	"database/sql"
	"log"
	"messagenow/exceptions"
)

type joinRoomRepositoryImpl struct {
	db *sql.DB
}

func NewJoinRoomRepository(db *sql.DB) JoinRoomRepository {
	return joinRoomRepositoryImpl{
		db: db,
	}
}

func (j joinRoomRepositoryImpl) Execute(ctx context.Context, roomID, userID int64) error {
	query := `
	INSERT INTO user_room (id_room, id_user) 
	VALUES (?, ?) ON DUPLICATE KEY UPDATE id_room = id_room`

	_, err := j.db.ExecContext(ctx, query, roomID, userID)
	if err != nil {
		log.Println("[joinRoomRepositoryImpl] Error ExecContext", err)
		return exceptions.NewInternalServerError(exceptions.InternalErrorMessage)
	}

	return nil
}
