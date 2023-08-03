package entities

type Message struct {
	Sender User   `json:"sender"`
	ID     int64  `json:"id"`
	Text   string `json:"text"`
}
