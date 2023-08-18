package repositories

import (
	"context"
	"database/sql"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
)

type getRoomsRepositoryImpl struct {
	db *sql.DB
}

func NewGetRoomsRepository(db *sql.DB) GetRoomsRepository {
	return getRoomsRepositoryImpl{
		db: db,
	}
}

func (g getRoomsRepositoryImpl) Execute(ctx context.Context, userID int64) ([]entities.Room, error) {
	//language=sql
	query := `
	SELECT r.id,
	       r.name,
	       r.image,
	       r.created_at,
	       r.modified_at
	FROM room r
	        INNER JOIN user_room ur ON r.id = ur.id_room AND
	                                   ur.id_user = ?
	        INNER JOIN user u1 ON ur.id_user = u1.id`

	rows, err := g.db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Println("[getRoomsRepositoryImpl] Error QueryContext", err)
		return nil, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}
	defer rows.Close()

	var rooms []entities.Room
	for rows.Next() {
		var room entities.Room
		err = rows.Scan(
			&room.ID,
			&room.Name,
			&room.ImageURL,
			&room.CreatedAt,
			&room.ModifiedAt,
		)
		if err != nil {
			log.Println("[getRoomsRepositoryImpl] Error Scan", err)
			return nil, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (g getRoomsRepositoryImpl) getRoomUsers(ctx context.Context, roomID int64) ([]entities.User, error) {
	//language=sql
	query := `
	SELECT u.id,
	       u.name,
	       u.image
	FROM user u
	    INNER JOIN user_room ur on u.id = ur.id_user AND
                               ur.id_room = ?`

	rows, err := g.db.QueryContext(ctx, query, roomID)
	if err != nil {
		log.Println("[getRoomUsers] Error QueryContext", err)
		return nil, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.ImageURL,
		)
		if err != nil {
			log.Println("[getRoomUsers] Error Scan", err)
			return nil, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
		}

		users = append(users, user)
	}

	return users, nil
}
