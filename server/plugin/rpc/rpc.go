package rpc

import (
	"context"
	"database/sql"
	"github.com/heroiclabs/nakama-common/runtime"
)

func getFirstWorld(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule) (string, error) {

	// List existing matches
	// that have been 1 & 4 players
	minSize := 1
	maxSize := 4
	matches, listErr := nk.MatchList(ctx, 1, false, "", &minSize, &maxSize, "") //local matches = nakama.match_list()

	// Return if listing error
	if listErr != nil {
		logger.Printf("Failed to list matches when grabing first world! Error: %v\n", listErr)
		return "", listErr
	}

	// If no matches exist, create one
	if len(matches) <= 0 {

		// Create match
		//params := map[string]interface{}{}
		matchID, createErr := nk.MatchCreate(ctx, "control", map[string]interface{}{})
		//return nakama.match_create("world_control", {})

		// Return if creation error
		if createErr != nil {
			logger.Printf("Failed to create match when grabing first world! Error: %v\n", createErr)
			return "", createErr
		}
		logger.Info("Successfully created new match!")

		// Return newly created match
		return matchID, nil

	} else {

		// Return first found match
		return matches[0].GetMatchId(), nil
	}

}

func GetWorldId(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	matchID, err := getFirstWorld(ctx, logger, nk)
	return matchID, err
}
