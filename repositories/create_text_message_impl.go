package repositories

import (
	"database/sql"
	"log"
	"messagenow/domain/entities"
)

type createTextMessageRepositoryImpl struct {
	db *sql.DB
}

func NewCreateTextMessageRepository(db *sql.DB) CreateTextMessageRepository {
	return createTextMessageRepositoryImpl{db: db}
}

func (c createTextMessageRepositoryImpl) Execute(Message entities.Message, senderID int64, addresseeID int64) error {
	tx, err := c.db.Begin()
	if err != nil {
		log.Println("[createTextMessageRepositoryImpl] Error Begin", err)
		return err
	}

	query := `
	INSERT INTO message (id_user, id_sender, id_addressee) 
	VALUES (?, ?, ?)`

	result, err := tx.Exec(query, senderID, senderID, addresseeID)
	if err != nil {
		_ = tx.Rollback()
		log.Println("[createTextMessageRepositoryImpl] Error ExecContext", err)
		return err
	}

	messageID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		log.Println("[createTextMessageRepositoryImpl] Error LastInsertId", err)
		return err
	}

	query = `
	INSERT INTO message_text (id_message, text) 
	VALUES (?, ?)`

	_, err = tx.Exec(query, messageID, Message.Text)
	if err != nil {
		_ = tx.Rollback()
		log.Println("[createTextMessageRepositoryImpl] Error ExecContext", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("[createTextMessageRepositoryImpl] Error Commit", err)
		return err
	}

	return nil
}
