package entities

type Message struct {
	SenderID    int64 `json:"senderID"`
	addresseeID int64 `json:"addresseeID"`
}

type MessageText struct {
	Text string `json:"text"`
}
