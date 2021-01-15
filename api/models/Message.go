package models

import "time"

// MessageType is the type of the message sent through Websockets
type MessageType int

const (
	// MessageTracking message to move player
	MessageTracking MessageType = 0
	// MessageMove message to move player
	MessageMove MessageType = 1
	// MessageRegister message to register player
	MessageRegister MessageType = 2
	// MessageUnregister to unregister player
	MessageUnregister MessageType = 3
	// MessageGetPlayers to get all other players
	MessageGetPlayers MessageType = 4
	// MessageCollide to test collision
	MessageCollide MessageType = 5
)

// Message  sent through Websockets
type Message struct {
	Type       MessageType        `json:"type"`
	Player     *Player            `json:"player_info"`
	ReceivedAt time.Time          `json:"received_at,omitempty"`
	SentAt     time.Time          `json:"sent_at,omitempty"`
	Players    *map[string]Player `json:"players,omitempty"`
}
