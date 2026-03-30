package main

import (
	"context"
	"database/sql"

	"tic-tac-toe/api"
	"tic-tac-toe/match"

	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	initializer runtime.Initializer,
) error {

	logger.Info("Initializing module")

	// Register Match
	if err := initializer.RegisterMatch(
		"tic_tac_toe",
		func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
			return &match.Match{}, nil
		},
	); err != nil {
		return err
	}

	// Register RPC
	if err := initializer.RegisterRpc("create_match", api.RpcCreateMatch); err != nil {
		return err
	}

	return nil
}
