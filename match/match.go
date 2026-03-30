package match

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

// Match struct (REQUIRED by Nakama)
type Match struct{}

/*
---------------- MATCH INIT ----------------
Called when match is created
*/
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

	return state, 2, "" // 2 ticks/sec
}

/*
---------------- JOIN ATTEMPT ----------------
Allow or reject player
*/
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

	return state, true, ""
}

/*
---------------- JOIN ----------------
Player enters match
*/
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
		s.Players = append(s.Players, p.GetUserId())
	}

	if len(s.Players) == 1 {
		s.Turn = s.Players[0]
		s.Status = "waiting"
	}

	if len(s.Players) == 2 {
		s.Status = "playing"
	}

	return s
}

/*
---------------- GAME LOOP ----------------
Handles moves
*/
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

		var data map[string]int
		if err := json.Unmarshal(msg.GetData(), &data); err != nil {
			continue
		}

		row := data["row"]
		col := data["col"]
		player := msg.GetUserId()

		// turn validation
		if player != s.Turn {
			continue
		}

		// already filled
		if s.Board[row][col] != "" {
			continue
		}

		// assign symbol
		symbol := "X"
		if len(s.Players) > 1 && player != s.Players[0] {
			symbol = "O"
		}

		s.Board[row][col] = symbol
		s.MovesCount++

		// switch turn
		if len(s.Players) == 2 {
			if player == s.Players[0] {
				s.Turn = s.Players[1]
			} else {
				s.Turn = s.Players[0]
			}
		}

		stateBytes, _ := json.Marshal(s)

		dispatcher.BroadcastMessage(1, stateBytes, nil, nil, true)
	}

	return s
}

/*
---------------- LEAVE ----------------
*/
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
	return state
}

/*
---------------- TERMINATE ----------------
*/
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

/*
---------------- SIGNAL ----------------
*/
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
