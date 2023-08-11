package views

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
	"messagenow/usecases"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type messageHttpModule struct {
	createMessageUseCase       usecases.CreateMessageUseCase
	getPreviousMessagesUseCase usecases.GetMessagesUseCase
	createRoomUseCase          usecases.CreateRoomUseCase
	joinRoomUseCase            usecases.JoinRoomUseCase
	getRoomsUseCase            usecases.GetRoomsUseCase
}

func NewMessageHTTPModule(
	createTextMessageUseCase usecases.CreateMessageUseCase,
	getPreviousMessagesUseCase usecases.GetMessagesUseCase,
	createRoomUseCase usecases.CreateRoomUseCase,
	joinRoomUseCase usecases.JoinRoomUseCase,
	getRoomsUseCase usecases.GetRoomsUseCase,
) ModuleHTTP {
	return messageHttpModule{
		createMessageUseCase:       createTextMessageUseCase,
		getPreviousMessagesUseCase: getPreviousMessagesUseCase,
		createRoomUseCase:          createRoomUseCase,
		joinRoomUseCase:            joinRoomUseCase,
		getRoomsUseCase:            getRoomsUseCase,
	}
}

func (m messageHttpModule) Setup(router *mux.Router) {
	router.HandleFunc("/rooms", m.createRoom).Methods(http.MethodPost)
	router.HandleFunc("/rooms", m.getRooms).Methods(http.MethodGet)
	router.HandleFunc("/rooms/{roomID}/join", m.joinRoom).Methods(http.MethodPost)
	router.HandleFunc("/rooms/{roomID}/messages", m.getPreviousMessages).Methods(http.MethodGet)
	router.HandleFunc("/rooms/{roomID}/ws", m.handleWebSocket)
	go m.broadCastMessages()
}

func (m messageHttpModule) createRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := m.createRoomUseCase.Execute(r.Context())
	if err != nil {
		log.Println("[getPreviousMessages] Error Execute", err)
		exceptions.HandleError(w, err)
		return
	}

	_, err = w.Write([]byte(strconv.FormatInt(roomID, 10)))
	if err != nil {
		log.Println("[getPreviousMessages] Error Write", err)
		exceptions.HandleError(w, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage))
		return
	}
}

func (m messageHttpModule) getRooms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(entities.User)
	userRooms, err := m.getRoomsUseCase.Execute(ctx, user.ID)
	if err != nil {
		log.Println("[getRooms] Error Execute", err)
		exceptions.HandleError(w, err)
		return
	}

	b, err := json.Marshal(userRooms)
	if err != nil {
		log.Println("[getRooms] Error Marshal", err)
		exceptions.HandleError(w, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("[getRooms] Error Write", err)
		exceptions.HandleError(w, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage))
		return
	}
}

func (m messageHttpModule) joinRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.ParseInt(mux.Vars(r)["roomID"], 10, 64)
	if err != nil {
		log.Println("[joinRoom] Error ParseInt", err)
		exceptions.HandleError(w, exceptions.NewBadRequestError("roomID is not valid"))
		return
	}

	ctx := r.Context()
	user := ctx.Value("user").(entities.User)
	err = m.joinRoomUseCase.Execute(ctx, roomID, user.ID)
	if err != nil {
		log.Println("[joinRoom] Error Execute", err)
		exceptions.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (m messageHttpModule) getPreviousMessages(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.ParseInt(mux.Vars(r)["roomID"], 10, 64)
	if err != nil {
		log.Println("[getPreviousMessages] Error ParseInt", err)
		exceptions.HandleError(w, exceptions.NewBadRequestError("roomID is not valid"))
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
		exceptions.HandleError(w, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("[getPreviousMessages] Error Write", err)
		exceptions.HandleError(w, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage))
		return
	}
}

type Client struct {
	entities.User
	conn *websocket.Conn
	mu   sync.Mutex
}

type ClientMessage struct {
	RoomID    int64         `json:"roomID"`
	User      entities.User `json:"sender"`
	Text      string        `json:"text"`
	CreatedAt time.Time     `json:"createdAt"`
}

var (
	rooms     = make(map[int64]map[*Client]bool)
	broadcast = make(chan ClientMessage)
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

	log.Println("userID", user.ID)

	roomID, err := strconv.ParseInt(mux.Vars(r)["roomID"], 10, 64)
	if err != nil {
		log.Println("[handleWebSocket] Error ParseInt", err)
		exceptions.HandleError(w, exceptions.NewBadRequestError("roomID is not valid"))
		return
	}

	client := &Client{conn: conn, User: user}

	m.addClient(roomID, client)

	fmt.Println("roomID", roomID)
	for roomClient, _ := range rooms[roomID] {
		fmt.Printf("%+v\n", roomClient.User.Name)
	}

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("[handleWebSocket] Error ReadJSON", err)
			break
		}

		broadcast <- ClientMessage{
			User:   client.User,
			RoomID: roomID,
			Text:   string(message),
		}
	}

	m.removeClient(roomID, client)
}

func (m messageHttpModule) addClient(id int64, client *Client) {
	client.mu.Lock()

	// create new map of clients in case of the room do not have any client.
	clients, ok := rooms[id]
	if !ok {
		clients = make(map[*Client]bool)
	}

	// if the client is not on the room clients map add it
	_, ok = clients[client]
	if !ok {
		clients[client] = true
	}

	rooms[id] = clients

	client.mu.Unlock()
}

func (m messageHttpModule) removeClient(id int64, client *Client) {
	client.mu.Lock()
	clients := rooms[id]
	delete(clients, client)
	rooms[id] = clients
	client.mu.Unlock()
}

func (m messageHttpModule) broadCastMessages() {
	for {
		clientMessage := <-broadcast
		fmt.Println("\n\n-----------------------------")
		err := m.createMessageUseCase.Execute(
			clientMessage.User.ID,
			clientMessage.RoomID,
			entities.Message{Text: clientMessage.Text},
		)
		if err != nil {
			log.Println("[broadCastMessages] Error createMessageUseCase.Execute", err)
		}

		for addressee := range rooms[clientMessage.RoomID] {
			addressee.mu.Lock()
			clientMessage.CreatedAt = time.Now()

			fmt.Printf("%+v\n", addressee.User.Name)

			err = addressee.conn.WriteJSON(clientMessage)
			addressee.mu.Unlock()
			if err != nil {
				log.Println("Error broadcasting message:", err)
			}
		}
	}
}
