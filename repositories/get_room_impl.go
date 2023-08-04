package repositories

import (
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

func (g getRoomRepositoryImpl) Execute(roomID int64, senderID int64) (*entities.Room, error) {
	//language=sql
	query := `
	SELECT id
	FROM room r
		INNER JOIN user_room ur1 on ur1.id_room = r.id AND 
								   ur1.id_user = ?
   WHERE r.id = ?`

	var room entities.Room
	err := g.db.QueryRow(query, roomID, senderID).Scan(&room.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("[getRoomRepositoryImpl] Error Scan", err)
		return nil, err
	}

	return &room, nil
}
