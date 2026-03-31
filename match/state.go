package match

// GameState stores the complete match state on the server.
type GameState struct {
	Board      [3][3]string `json:"board"`
	Players    []string     `json:"players"`
	Turn       string       `json:"turn,omitempty"`
	Winner     string       `json:"winner,omitempty"`
	MovesCount int          `json:"moves_count"`
	Status     string       `json:"status,omitempty"` // waiting, playing, finished
}
