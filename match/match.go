package match

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

type Match struct{}

// MatchInit runs when a new authoritative match is created.
func (m *Match) MatchInit(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	params map[string]interface{},
) (interface{}, int, string) {

	state := &GameState{
		Board:   [3][3]string{},
		Players: []string{},
		Status:  "waiting",
	}

	logger.Info("Match initialized")

	return state, 2, ""
}

// MatchJoinAttempt validates whether a player can join.
func (m *Match) MatchJoinAttempt(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher,
	tick int64,
	state interface{},
	presence runtime.Presence,
	metadata map[string]string,
) (interface{}, bool, string) {

	s := state.(*GameState)

	if len(s.Players) >= 2 {
		return s, false, "match is full"
	}

	return s, true, ""
}

// MatchJoin adds player(s) to the match and broadcasts updated state.
func (m *Match) MatchJoin(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher,
	tick int64,
	state interface{},
	presences []runtime.Presence,
) interface{} {

	s := state.(*GameState)

	for _, p := range presences {
		exists := false
		for _, playerID := range s.Players {
			if playerID == p.GetUserId() {
				exists = true
				break
			}
		}

		if !exists {
			s.Players = append(s.Players, p.GetUserId())
		}
	}

	if len(s.Players) == 1 {
		s.Turn = s.Players[0]
		s.Status = "waiting"
	}

	if len(s.Players) == 2 {
		s.Turn = s.Players[0]
		s.Status = "playing"
	}

	stateBytes, err := json.Marshal(s)
	if err != nil {
		logger.Error("Failed to serialize state on join: %v", err)
		return s
	}

	_ = dispatcher.BroadcastMessage(1, stateBytes, nil, nil, true)

	logger.Info("Player joined. Players count: %d", len(s.Players))

	return s
}

// MatchLoop processes moves sent by players.
func (m *Match) MatchLoop(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher,
	tick int64,
	state interface{},
	messages []runtime.MatchData,
) interface{} {

	s := state.(*GameState)

	for _, msg := range messages {
		if s.Status != "playing" {
			continue
		}

		if s.Winner != "" {
			continue
		}

		var data map[string]int
		if err := json.Unmarshal(msg.GetData(), &data); err != nil {
			logger.Error("Invalid move data: %v", err)
			continue
		}

		row := data["row"]
		col := data["col"]
		player := msg.GetUserId()

		if row < 0 || row > 2 || col < 0 || col > 2 {
			logger.Warn("Invalid position")
			continue
		}

		if player != s.Turn {
			logger.Warn("Not player's turn")
			continue
		}

		if s.Board[row][col] != "" {
			logger.Warn("Cell already occupied")
			continue
		}

		symbol := "X"
		if len(s.Players) > 1 && player == s.Players[1] {
			symbol = "O"
		}

		s.Board[row][col] = symbol
		s.MovesCount++

		winnerSymbol := CheckWinner(s.Board)
		if winnerSymbol != "" {
			if winnerSymbol == "X" {
				s.Winner = s.Players[0]
			} else if len(s.Players) > 1 {
				s.Winner = s.Players[1]
			}
			s.Status = "finished"
		} else if s.MovesCount == 9 {
			s.Status = "finished"
		} else {
			if len(s.Players) == 2 {
				if player == s.Players[0] {
					s.Turn = s.Players[1]
				} else {
					s.Turn = s.Players[0]
				}
			}
		}

		stateBytes, err := json.Marshal(s)
		if err != nil {
			logger.Error("Failed to serialize state: %v", err)
			continue
		}

		_ = dispatcher.BroadcastMessage(1, stateBytes, nil, nil, true)
	}

	return s
}

// MatchLeave handles player leaving.
func (m *Match) MatchLeave(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher,
	tick int64,
	state interface{},
	presences []runtime.Presence,
) interface{} {

	s := state.(*GameState)

	for _, p := range presences {
		filtered := make([]string, 0, len(s.Players))
		for _, playerID := range s.Players {
			if playerID != p.GetUserId() {
				filtered = append(filtered, playerID)
			}
		}
		s.Players = filtered
	}

	if len(s.Players) == 0 {
		s.Status = "finished"
	} else if len(s.Players) == 1 {
		s.Status = "waiting"
		s.Turn = s.Players[0]
	}

	stateBytes, err := json.Marshal(s)
	if err == nil {
		_ = dispatcher.BroadcastMessage(1, stateBytes, nil, nil, true)
	}

	return s
}

// MatchTerminate handles cleanup.
func (m *Match) MatchTerminate(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher,
	tick int64,
	state interface{},
	graceSeconds int,
) interface{} {
	return state
}

// MatchSignal handles external signals.
func (m *Match) MatchSignal(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher,
	tick int64,
	state interface{},
	data string,
) (interface{}, string) {
	return state, ""
}
