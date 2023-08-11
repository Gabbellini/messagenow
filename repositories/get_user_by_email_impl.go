package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
)

type getUserByEmailRepositoryImpl struct {
	db *sql.DB
}

func NewGetUserByEmailRepository(db *sql.DB) GetUserByEmailRepository {
	return getUserByEmailRepositoryImpl{db: db}
}

func (g getUserByEmailRepositoryImpl) Execute(ctx context.Context, email string) (*entities.User, error) {
	//language=sql
	query := `
	SELECT id, 
	       name, 
	       image, 
	       email
	FROM user
	WHERE email = ?`

	var user entities.User
	err := g.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.ImageURL,
		&user.Credential.Email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		log.Println("[getUserByEmailRepositoryImpl] Error Scan", err)
		return nil, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}

	return &user, err
}
