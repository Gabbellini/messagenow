package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
	"messagenow/usecases"
	"net/http"
	"strconv"
	"sync"
)

type messageHttpModule struct {
	createTextMessageUseCase   usecases.CreateTextMessageUseCase
	getPreviousMessagesUseCase usecases.GetPreviousMessagesUseCase
}

func NewMessageHTTPModule(
	createTextMessageUseCase usecases.CreateTextMessageUseCase,
	getPreviousMessagesUseCase usecases.GetPreviousMessagesUseCase,
) ModuleHTTP {
	return messageHttpModule{
		createTextMessageUseCase:   createTextMessageUseCase,
		getPreviousMessagesUseCase: getPreviousMessagesUseCase,
	}
}

func (m messageHttpModule) Setup(router *mux.Router) {
	router.HandleFunc("/ws", m.handleWebSocket)
	router.HandleFunc("/messages/{roomID}", m.getPreviousMessages)
	go m.broadCastMessages()
}

type Client struct {
	entities.User
	conn *websocket.Conn
	mu   sync.Mutex
}

type ClientTextMessage struct {
	ClientID   int64  `json:"id"`
	ClientName string `json:"name"`
	Text       string `json:"text"`
}

var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan ClientTextMessage)
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (m messageHttpModule) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[handleWebSocket] Error Upgrade", err)
		return
	}

	ctx := r.Context()
	user := ctx.Value("user").(entities.User)

	client := &Client{conn: conn, User: user}

	m.addClient(client)
	defer m.removeClient(client)

	log.Println("client Connected", client.Name)

	var message entities.Message
	for {
		// Wait for the JSON message from the client
		err = client.conn.ReadJSON(&message)
		if err != nil {
			log.Println("[handleWebSocket] Error ReadJSON", err)
			return
		}

		// handle Message
		m.handleMessage(client, message)
	}
}

func (m messageHttpModule) addClient(client *Client) {
	client.mu.Lock()
	clients[client] = true
	client.mu.Unlock()
}

func (m messageHttpModule) removeClient(client *Client) {
	client.mu.Lock()
	delete(clients, client)
	client.mu.Unlock()
}

func (m messageHttpModule) handleMessage(sender *Client, message entities.Message) {
	// Process the received message here (e.g., handle commands, etc.)
	// For simplicity, we just broadcast the message to all connected clients.

	broadcast <- ClientTextMessage{
		ClientID:   sender.ID,
		ClientName: sender.Name,
		Text:       message.Text,
	}
}

func (m messageHttpModule) broadCastMessages() {
	for {
		clientMessage := <-broadcast
		for addressee := range clients {
			addressee.mu.Lock()
			err := addressee.conn.WriteJSON(clientMessage)
			addressee.mu.Unlock()
			if err != nil {
				log.Println("Error broadcasting message:", err)
			}
			err = m.createTextMessageUseCase.Execute(entities.Message{Text: clientMessage.Text}, clientMessage.ClientID, addressee.ID)
			if err != nil {
				log.Println("[broadCastMessages] Error createTextMessageUseCase.Execute", err)
			}
		}
	}
}

func (m messageHttpModule) getPreviousMessages(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.ParseInt(mux.Vars(r)["roomID"], 64, 10)
	if err != nil {
		log.Println("[getPreviousMessages] Error ParseInt", err)
		exceptions.HandleError(w, exceptions.NewBadRequestError("clientID is not valid"))
		return
	}

	ctx := r.Context()
	user := ctx.Value("user").(entities.User)
	messages, err := m.getPreviousMessagesUseCase.Execute(ctx, user.ID, roomID)
	if err != nil {
		log.Println("[getPreviousMessages] Error Execute", err)
		exceptions.HandleError(w, err)
		return
	}

	b, err := json.Marshal(messages)
	if err != nil {
		log.Println("[getPreviousMessages] Error Marshal", err)
		exceptions.HandleError(w, exceptions.NewInternalServerError(exceptions.InternalErrorMessage))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("[getPreviousMessages] Error Write", err)
		exceptions.HandleError(w, exceptions.NewInternalServerError(exceptions.InternalErrorMessage))
		return
	}
}
