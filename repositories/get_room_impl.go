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

func (g getRoomRepositoryImpl) Execute(roomID int64, addresseeID int64) (*entities.Room, error) {
	//language=sql
	query := `
	SELECT id
	FROM room 
		INNER JOIN user_room ur1 on ur1.id_room = room.id AND 
								   ur1.id_user = ?
		INNER JOIN user_room ur2 on ur2.id_room = room.id AND 
						   ur2.id_user = ?`

	var room entities.Room
	err := g.db.QueryRow(query, roomID, addresseeID).Scan(&room.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("[getRoomRepositoryImpl] Error Scan", err)
		return nil, err
	}

	return &room, nil
}
