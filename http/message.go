package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"messagenow/domain/entities"
	"messagenow/usecases"
	"net/http"
)

type messageHttpModule struct {
	createTextMessage usecases.CreateTextMessageUseCase
}

func NewMessageHTTPModule(createTextMessage usecases.CreateTextMessageUseCase) ModuleHTTP {
	return messageHttpModule{createTextMessage: createTextMessage}
}

func (m messageHttpModule) Setup(router *mux.Router) {
	router.HandleFunc("/message/{userID}", m.messageHandler)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow any origin for the WebSocket connection (for demo purposes)
		return true
	},
}

func (m messageHttpModule) messageHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	for {
		// Read message from the WebSocket client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Print the received message to the console
		fmt.Printf("Received message: %s\n", message)

		userID := int64(4)
		senderID := userID
		addresseeID := int64(5)

		err = m.createTextMessage.Execute(r.Context(), entities.MessageText{Text: string(message)}, senderID, addresseeID)
		if err != nil {
			log.Println("[messageHandler] Error createTextMessage.Execute", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond back to the WebSocket client with the same message
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
