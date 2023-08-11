package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
)

type getUserByIDRepositoryImpl struct {
	db *sql.DB
}

func NewGetUserByIDRepository(db *sql.DB) GetUserByIDRepository {
	return getUserByIDRepositoryImpl{db: db}
}

func (g getUserByIDRepositoryImpl) Execute(ctx context.Context, userID int64) (*entities.User, error) {
	//language=sql
	query := `
	SELECT id, 
	       name, 
	       image, 
	       email
	FROM user
	WHERE id = ?`

	var user entities.User
	err := g.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.ImageURL,
		&user.Credential.Email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		log.Println("[getUserByIDRepositoryImpl] Error Scan", err)
		return nil, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}

	return &user, err
}
