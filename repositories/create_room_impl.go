package repositories

import (
	"context"
	"database/sql"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
)

type createRoomRepositoryImpl struct {
	db *sql.DB
}

func NewCreateRoomRepository(db *sql.DB) CreateRoomRepository {
	return createRoomRepositoryImpl{
		db: db,
	}
}

func (c createRoomRepositoryImpl) Execute(ctx context.Context, room entities.Room) (int64, error) {
	//language=sql
	query := `
	INSERT INTO room (type, image) VALUES (?, ?)`

	result, err := c.db.ExecContext(ctx, query, room.Type, room.ImageURL)
	if err != nil {
		log.Println("[createRoom] Error ExecContext", err)
		return 0, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("[createRoom] Error LastInsertId", err)
		return 0, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}

	return id, nil
}
