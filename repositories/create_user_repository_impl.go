package repositories

import (
	"context"
	"database/sql"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
)

type createUserRepositoryImpl struct {
	db *sql.DB
}

func NewCreateUserRepository(db *sql.DB) CreateUserRepository {
	return createUserRepositoryImpl{db: db}
}

func (c createUserRepositoryImpl) Execute(ctx context.Context, user entities.User) (int64, error) {
	query := `
	INSERT INTO user (name, image, email, password) 
	VALUES (?, ?, ?, ?)`

	result, err := c.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.ImageURL,
		user.Credential.Email,
		user.Credential.Password,
	)
	if err != nil {
		log.Println("[createUserRepositoryImpl] Error ExecContext", err)
		return 0, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("[createUserRepositoryImpl] Error lastInsertedId", err)
		return 0, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage)
	}

	return id, nil
}
