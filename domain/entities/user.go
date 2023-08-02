package entities

type User struct {
	Name     string  `json:"name"`
	ImageURL *string `json:"image"`
}
