package repositories

import (
	"context"
	"database/sql"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
)

type getRoomUsersRepositoryImpl struct {
	db *sql.DB
}

func NewGetRoomUsersRepository(db *sql.DB) GetRoomUsersRepository {
	return getRoomUsersRepositoryImpl{
		db: db,
	}
}

func (g getRoomUsersRepositoryImpl) Execute(ctx context.Context, roomID int64) ([]entities.User, error) {
	//language=sql
	query := `
	SELECT u.id,
	       u.name,
	       u.image
	FROM user u
	INNER JOIN user_room ur ON u.id = ur.id_user AND
	                           ur.id_room = ?`

	rows, err := g.db.QueryContext(ctx, query, roomID)
	if err != nil {
		log.Println("[getRoomUsersRepositoryImpl] Error QueryContext", err)
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
			log.Println("[getRoomUsersRepositoryImpl] Error QueryContext", err)
			return nil, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
		}

		users = append(users, user)
	}

	return users, nil
}
