package entities

type MessageUser struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Image *string `json:"image"`
}

type Message struct {
	Sender    MessageUser `json:"sender"`
	Addressee MessageUser `json:"-"`
	ID        int64       `json:"id"`
	Text      string      `json:"text"`
}
