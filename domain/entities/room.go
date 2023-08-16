package entities

import "time"

const RoomTypeSingle = 1
const RoomTypeChat = 2
const RoomTypeGroup = 3

const RoomStatusOk = 1
const RoomStatusMissingUsers = 2

type Room struct {
	ID         int64     `json:"id"`
	Type       int64     `json:"type"`
	Name       string    `json:"name"`
	ImageURL   *string   `json:"imageURL"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}
