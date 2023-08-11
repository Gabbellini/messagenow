package repositories

import (
	"context"
	"database/sql"
	"log"
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

func (c createRoomRepositoryImpl) Execute(ctx context.Context) (int64, error) {
	query := `
	INSERT INTO room (created_at) VALUES (CURRENT_TIMESTAMP)`

	result, err := c.db.ExecContext(ctx, query)
	if err != nil {
		log.Println("[createRoomRepositoryImpl] createRoom Error Exec", err)
		return 0, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("[createRoomRepositoryImpl] createRoom Error LastInsertId", err)
		return 0, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}

	return id, nil
}
