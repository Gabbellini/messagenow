package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
)

type getUserRepositoryImpl struct {
	db *sql.DB
}

func NewGetUserRepository(db *sql.DB) GetUserRepository {
	return getUserRepositoryImpl{db: db}
}

func (g getUserRepositoryImpl) Execute(ctx context.Context, userID int64) (*entities.User, error) {
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
			log.Println("[loginRepositoryImpl] Error sql.ErrNoRows", err)
			return nil, exceptions.NewBadRequestError("Usuário não encontrado")
		}

		log.Println("[loginRepositoryImpl] Error Scan", err)
		return nil, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}

	return &user, err
}
