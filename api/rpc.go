package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/heroiclabs/nakama-common/runtime"
)

type CreateMatchResponse struct {
	MatchID string `json:"match_id"`
}

// RPC: create match
func RpcCreateMatch(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB, // ✅ FIXED
	nk runtime.NakamaModule,
	payload string,
) (string, error) {

	matchID, err := nk.MatchCreate(ctx, "tic_tac_toe", nil)
	if err != nil {
		return "", err
	}

	res := CreateMatchResponse{MatchID: matchID}
	bytes, _ := json.Marshal(res)

	return string(bytes), nil
}
