package match

// GameState = full game data stored on server
type GameState struct {
	Board      [3][3]string `json:"board"`       // 3x3 grid
	Players    []string     `json:"players"`     // player IDs
	Turn       string       `json:"turn"`        // whose turn
	Winner     string       `json:"winner"`      // winner ID
	MovesCount int          `json:"moves_count"` // total moves
	Status     string       `json:"status"`      // waiting, playing, finished
}
