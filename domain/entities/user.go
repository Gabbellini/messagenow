package entities

type User struct {
	ID         int64       `json:"id"`
	Name       string      `json:"name"`
	ImageURL   *string     `json:"image"`
	Credential Credentials `json:"-"`
}
