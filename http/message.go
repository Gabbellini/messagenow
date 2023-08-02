package http

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"messagenow/domain/entities"
	"messagenow/usecases"
	"net/http"
	"sync"
)

type messageHttpModule struct {
	createTextMessage usecases.CreateTextMessageUseCase
}

func NewMessageHTTPModule(createTextMessage usecases.CreateTextMessageUseCase) ModuleHTTP {
	return messageHttpModule{createTextMessage: createTextMessage}
}

func (m messageHttpModule) Setup(router *mux.Router) {
	router.HandleFunc("/ws", m.handleWebSocket)
	go m.broadCastMessages()
}

type Client struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan entities.MessageText)
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

	client := &Client{conn: conn}
	m.addClient(client)
	defer m.removeClient(client)

	var message entities.MessageText
	for {
		err = client.conn.ReadJSON(&message)
		if err != nil {
			log.Println("[handleWebSocket] Error ReadJSON", err)
			return
		}

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

func (m messageHttpModule) handleMessage(sender *Client, message entities.MessageText) {
	// Process the received message here (e.g., handle commands, etc.)
	// For simplicity, we just broadcast the message to all connected clients.

	broadcast <- message
}

func (m messageHttpModule) broadCastMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			client.mu.Lock()
			err := client.conn.WriteJSON(msg)
			client.mu.Unlock()
			if err != nil {
				log.Println("Error broadcasting message:", err)
			}
		}
	}
}
