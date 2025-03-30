package entity

// --- Game State Types ---

// Position represents a coordinate on the game grid.
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Direction represents the movement vector of a snake.
type Direction struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Predefined directions
var (
	DirUp    = Direction{X: 0, Y: -1}
	DirDown  = Direction{X: 0, Y: 1}
	DirLeft  = Direction{X: -1, Y: 0}
	DirRight = Direction{X: 1, Y: 0}
	DirNone  = Direction{X: 0, Y: 0} // Optional: Represents no movement/initial state
)

// GameState represents the current state of the game session.
type GameState string

const (
	StateWaiting  GameState = "waiting"
	StateReady    GameState = "ready"
	StatePlaying  GameState = "playing"
	StateGameOver GameState = "gameover"
)

// PlayerState represents the state of a single player in the game.
// This is often sent from the client to the server.
type PlayerState struct {
	ID         string     `json:"id"`         // User ID (Should match the user sending the state)
	Snake      []Position `json:"snake"`      // List of positions forming the snake's body
	Direction  Direction  `json:"direction"`  // Current direction of movement
	Food       Position   `json:"food"`       // Current food position (as perceived by this client?)
	Score      int        `json:"score"`      // Player's score
	GameState  GameState  `json:"gameState"`  // Current game state (e.g., 'playing')
	LastUpdate int64      `json:"lastUpdate"` // Timestamp of the last state update (client-side)
	Joined     bool       `json:"joined"`     // Whether the player has successfully joined (less relevant in state updates?)
}

// GameWinner represents the winner of a game. Can be a user ID, 'draw', or empty.
type GameWinner string

const (
	WinnerDraw GameWinner = "draw"
)

// Game Socket Events
const (
	GameEventReply       = "reply" // General reply event
	GameEventError       = "error" // General error event
	GameEventCreateRoom  = "create_room"
	GameEventJoinRoom    = "join_room"
	GameEventLeaveRoom   = "leave_room"
	GameEventStartGame   = "start_game"
	GameEventEndGame     = "end_game"
	GameEventMove        = "move"
	GameEventPlayerState = "player_state" // For broadcasting full player states
	GameEventChat        = "chat"
	GameEventReconnect   = "reconnect"
	GameEventDisconnect  = "disconnect"
)
