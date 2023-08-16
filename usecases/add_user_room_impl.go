package usecases

import (
	"context"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
	"messagenow/repositories"
)

type addUserRoomUseCaseImpl struct {
	getRoomRepository       repositories.GetRoomRepository
	getRoomUsersRepository  repositories.GetRoomUsersRepository
	joinRoomRepository      repositories.JoinRoomRepository
	setRoomStatusRepository repositories.SetRoomStatusRepository
}

func NewAddUserRoomUseCase(
	joinRoomRepository repositories.JoinRoomRepository,
	getRoomRepository repositories.GetRoomRepository,
	getRoomUsersRepository repositories.GetRoomUsersRepository,
	setRoomStatusRepository repositories.SetRoomStatusRepository,
) AddUserRoomUseCase {
	return addUserRoomUseCaseImpl{
		joinRoomRepository:      joinRoomRepository,
		getRoomRepository:       getRoomRepository,
		getRoomUsersRepository:  getRoomUsersRepository,
		setRoomStatusRepository: setRoomStatusRepository,
	}
}

func (a addUserRoomUseCaseImpl) Execute(ctx context.Context, user entities.User, roomID, userID int64) error {
	room, err := a.getRoomRepository.Execute(ctx, roomID, user.ID)
	if err != nil {
		log.Println("[addUserRoomUseCaseImpl] Error getRoomRepository", err)
		return err
	}

	users, err := a.getRoomUsersRepository.Execute(ctx, roomID)
	if err != nil {
		log.Println("[addUserRoomUseCaseImpl] Error getRoomUsersRepository", err)
		return err
	}

	numberOfUsersBeforeAdd := len(users)
	if room.Type == entities.RoomTypeChat && numberOfUsersBeforeAdd == 2 {
		return exceptions.NewBadRequestError("Está sala já está cheia")
	}

	err = a.joinRoomRepository.Execute(ctx, roomID, userID)
	if err != nil {
		log.Println("[addUserRoomUseCaseImpl] Error joinRoomRepository", err)
		return err
	}

	if numberOfUsersBeforeAdd == 1 {
		err = a.setRoomStatusRepository.Execute(ctx, roomID, entities.RoomStatusOk)
		if err != nil {
			log.Println("[addUserRoomUseCaseImpl] Error setRoomStatusRepository", err)
			return err
		}
	}

	return nil
}
