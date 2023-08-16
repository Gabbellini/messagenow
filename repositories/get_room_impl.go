package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"messagenow/domain/entities"
)

type getRoomRepositoryImpl struct {
	db *sql.DB
}

func NewGetRoomRepository(db *sql.DB) GetRoomRepository {
	return getRoomRepositoryImpl{
		db: db,
	}
}

func (g getRoomRepositoryImpl) Execute(ctx context.Context, roomID int64, userID int64) (*entities.Room, error) {
	//language=sql
	query := `
	SELECT r.id,
	       r.image,
	       r.created_at
	FROM room r
		INNER JOIN user_room ur on ur.id_room = r.id AND 
								   ur.id_user = ?
    WHERE r.id = ?`

	var room entities.Room
	err := g.db.QueryRowContext(ctx, query, roomID, userID).Scan(&room.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("[getRoomRepositoryImpl] Error Scan", err)
		return nil, err
	}

	return &room, nil
}
