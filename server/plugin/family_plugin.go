package main

import (
	"context"
	"database/sql"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/josephbmanley/family/server/plugin/control"
	"github.com/josephbmanley/family/server/plugin/rpc"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Loaded family plugin!")

	if err := initializer.RegisterMatch("control", func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return &control.Match{}, nil
	}); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("get_world_id", rpc.GetWorldId); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}
	return nil
}
