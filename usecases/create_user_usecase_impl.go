package usecases

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
	"messagenow/repositories"
	"net/mail"
	"strings"
)

type createUserUseCaseImpl struct {
	getUserByEmailRepository repositories.GetUserByEmailRepository
	createUserRepository     repositories.CreateUserRepository
}

func NewCreateUserUseCase(
	createUserRepository repositories.CreateUserRepository,
	getUserByEmailRepository repositories.GetUserByEmailRepository,
) CreateUserUseCase {
	return createUserUseCaseImpl{
		createUserRepository:     createUserRepository,
		getUserByEmailRepository: getUserByEmailRepository,
	}
}

func (c createUserUseCaseImpl) Execute(ctx context.Context, user entities.User) (int64, error) {
	err := c.processUser(&user)
	if err != nil {
		log.Println("[createUserUseCaseImpl] Error validateUser", err)
		return 0, err
	}

	emailOwner, err := c.getUserByEmailRepository.Execute(ctx, user.Credential.Email)
	if err != nil {
		log.Println("[createUserUseCaseImpl] Error getUserByEmailRepository", err)
		return 0, err
	}

	if emailOwner != nil {
		return 0, exceptions.NewBadRequestError("Este e-mail já está em uso")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Credential.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[createUserUseCaseImpl] Error GenerateFromPassword", err)
		return 0, err
	}
	user.Credential.Password = string(passwordHash)

	return c.createUserRepository.Execute(ctx, user)
}

func (c createUserUseCaseImpl) processUser(user *entities.User) error {
	if user.Name = strings.TrimSpace(user.Name); user.Name == "" {
		return exceptions.NewBadRequestError("Nome não pode ser vazio")
	}

	if user.Credential.Email = strings.TrimSpace(user.Credential.Email); user.Credential.Email == "" {
		return exceptions.NewBadRequestError("Email não pode ser vazio")
	}

	_, err := mail.ParseAddress(user.Credential.Email)
	if err != nil {
		return exceptions.NewBadRequestError("Email não é válido")
	}

	if user.Credential.Password = strings.TrimSpace(user.Credential.Password); user.Credential.Password == "" {
		return exceptions.NewBadRequestError("Senha não pode ser vazia")
	}

	return nil
}
