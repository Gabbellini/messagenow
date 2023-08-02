package entities

type Message struct {
	SenderID    int64 `json:"senderID"`
	AddresseeID int64 `json:"addresseeID"`
}

type MessageText struct {
	Message
	Text string `json:"text"`
}
