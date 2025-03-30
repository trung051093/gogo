package dto

import "gogo/modules/snake_battle/entity"

// CreateRoomReq is the request payload for creating a new game room.
// It might be empty or contain optional settings in the future.
type CreateRoomReq struct {
	Ref        string `json:"ref"`                  // Reference ID for client-side tracking
	MaxPlayers *int   `json:"maxPlayers,omitempty"` // Optional: Specify max players for the room
}

// CreateRoomRes is the response payload after successfully creating a room.
type CreateRoomRes struct {
	entity.Room
}

// JoinRoomReq is the request payload for joining an existing game room.
type JoinRoomReq struct {
	Ref      string `json:"ref"`      // Reference ID for client-side tracking
	RoomCode string `json:"roomCode"` // The code of the room to join
}

// LeaveRoomReq is the request payload for leaving a game room.
type LeaveRoomReq struct {
	Ref      string `json:"ref"`      // Reference ID for client-side tracking
	RoomCode string `json:"roomCode"` // Usually inferred from the connection's current room
}

// MoveReq is the request payload for sending player movement/state.
// For snake game, sending only the new direction might be sufficient.
// Alternatively, send the full player state for synchronization.
type MoveReq struct {
	Ref   string             `json:"ref"` // Reference ID for client-side tracking
	State entity.PlayerState `json:"playerState"`
}

// Reply is a generic reply structure for socket events.
type Reply struct {
	Ok    bool   `json:"ok"`              // Indicates if the operation was successful
	Ref   string `json:"ref"`             // Reference ID from the original request
	Error string `json:"error,omitempty"` // Error message if Ok is false
	Data  any    `json:"data,omitempty"`  // Optional data payload (can hold any type)
}
