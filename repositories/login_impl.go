package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
)

type loginRepositoryImpl struct {
	db *sql.DB
}

func NewLoginRepository(db *sql.DB) LoginRepository {
	return loginRepositoryImpl{db: db}
}

func (l loginRepositoryImpl) Execute(ctx context.Context, credential entities.Credentials) (*entities.User, error) {
	//language=sql
	query := `
	SELECT id, 
	       name, 
	       image,
	       email,
	       password
	FROM user
	WHERE email = ?`

	var user entities.User
	err := l.db.QueryRowContext(ctx, query, credential.Email).Scan(
		&user.ID,
		&user.Name,
		&user.ImageURL,
		&user.Credential.Email,
		&user.Credential.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[loginRepositoryImpl] Error sql.ErrNoRows", err)
			return nil, exceptions.NewForbiddenError(exceptions.ForbiddenMessage)
		}

		log.Println("[loginRepositoryImpl] Error Scan", err)
		return nil, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}

	return &user, err
}
