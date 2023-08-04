package entities

import "time"

type Room struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Users      []User    `json:"users"`
}
