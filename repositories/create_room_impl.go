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

func (c createRoomRepositoryImpl) Execute(ctx context.Context) (*entities.Room, error) {
	query := `
	INSERT INTO room (created_at) VALUES (CURRENT_TIMESTAMP)`

	result, err := c.db.ExecContext(ctx, query)
	if err != nil {
		log.Println("[createRoomRepositoryImpl] createRoom Error Exec", err)
		return nil, exceptions.NewInternalServerError(exceptions.InternalErrorMessage)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("[createRoomRepositoryImpl] createRoom Error LastInsertId", err)
		return nil, exceptions.NewInternalServerError(exceptions.InternalErrorMessage)
	}

	return &entities.Room{ID: id}, nil
}
