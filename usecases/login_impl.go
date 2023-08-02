package usecases

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
	"messagenow/repositories"
	"strings"
)

type loginUseCaseImpl struct {
	loginRepository repositories.LoginRepository
}

func NewLoginUseCase(loginRepository repositories.LoginRepository) LoginUseCase {
	return loginUseCaseImpl{loginRepository: loginRepository}
}

func (l loginUseCaseImpl) Execute(ctx context.Context, credential entities.Credential) (*entities.User, error) {
	if credential.Email = strings.TrimSpace(credential.Email); credential.Email == "" {
		return nil, exceptions.NewForbiddenError(exceptions.ForbiddenMessage)
	}

	if credential.Password = strings.TrimSpace(credential.Password); credential.Password == "" {
		return nil, exceptions.NewForbiddenError(exceptions.ForbiddenMessage)
	}

	user, err := l.loginRepository.Execute(ctx, credential)
	if err != nil {
		log.Println("[Login] Error Execute", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Credential.Password), []byte(credential.Password))
	if err != nil {
		log.Println("[Login] Error CompareHashAndPassword", err)
		return nil, exceptions.NewForbiddenError(exceptions.ForbiddenMessage)
	}

	return user, nil
}
