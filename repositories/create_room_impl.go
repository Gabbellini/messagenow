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

func (c createRoomRepositoryImpl) Execute(ctx context.Context, roomID int64, addresseeID int64) (*entities.Room, error) {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("[createRoomRepositoryImpl] Error Begin", err)
		return nil, exceptions.NewInternalServerError(exceptions.InternalErrorMessage)
	}

	room, err := c.createRoom(ctx, tx)
	if err != nil {
		_ = tx.Rollback()
		log.Println("[createRoomRepositoryImpl] Error createRoom", err)
		return nil, err
	}

	err = c.createUsersOnRoom(ctx, tx, room.ID, roomID, addresseeID)
	if err != nil {
		_ = tx.Rollback()
		log.Println("[createRoomRepositoryImpl] Error createUsersOnRoom", err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		log.Println("[createRoomRepositoryImpl] Error Commit", err)
		return nil, exceptions.NewInternalServerError(exceptions.InternalErrorMessage)
	}

	return room, nil
}

func (c createRoomRepositoryImpl) createRoom(ctx context.Context, tx *sql.Tx) (*entities.Room, error) {
	query := `
	INSERT INTO room (created_at) VALUES (CURRENT_TIMESTAMP)`

	result, err := tx.ExecContext(ctx, query)
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

func (c createRoomRepositoryImpl) createUsersOnRoom(ctx context.Context, tx *sql.Tx, roomID, userID, addresseeID int64) error {
	//language=sql
	query := `
	INSERT INTO user_room (id_room, id_user) 
	VALUES (?, ?)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Println("[createRoomRepositoryImpl] Error", err)
		return exceptions.NewInternalServerError(exceptions.InternalErrorMessage)
	}

	users := []struct {
		ID int64
	}{
		{ID: userID},
		{ID: addresseeID},
	}

	for _, user := range users {
		_, err = stmt.Exec(roomID, user.ID)
		if err != nil {
			log.Println("[createRoomRepositoryImpl] Error stmt.Exec()", err)
			return exceptions.NewInternalServerError(exceptions.InternalErrorMessage)
		}
	}

	return nil
}
