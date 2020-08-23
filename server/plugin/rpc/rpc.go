package rpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/josephbmanley/family/server/plugin/entities"
	"github.com/josephbmanley/family/server/plugin/gameworld"
)

func getFirstWorld(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule) (string, error) {

	// List existing matches
	// that have been 1 & 32 players
	minSize := 1
	maxSize := 32

	// Lists server authorative servers
	if matches, err := nk.MatchList(ctx, 1, true, "", &minSize, &maxSize, ""); err != nil {
		logger.Printf("Failed to list matches when grabing first world! Error: %v\n", err)
		return "", err
	} else {
		// If no matches exist, create one
		if len(matches) <= 0 {

			// Create match
			matchID, createErr := nk.MatchCreate(ctx, "control", map[string]interface{}{})

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

}

func GetWorldId(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	matchID, err := getFirstWorld(ctx, logger, nk)
	return matchID, err
}

func CreateCharacter(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if ok {
		dataExist, err := entities.PlayerDataExists(ctx, nk, userID)
		if err != nil {
			logger.Error(err.Error())
			return err.Error(), err
		}

		if dataExist {
			return "Already Exists Exception", errors.New("user already has a character")
		} else {
			playerData := entities.PlayerSaveData{}
			err := json.Unmarshal([]byte(payload), &playerData)
			if err != nil {
				logger.Error("Failed to load data from client: %s", err.Error())
				return "Failed to load data from client!", err
			}
			player := entities.PlayerEntity{
				Name:    playerData.Name,
				Faction: gameworld.Faction(playerData.Faction),
			}
			saveErr := player.SaveUserID(ctx, nk, userID)
			if saveErr != nil {
				logger.Error("Failed to write data to storage on create: %s", err.Error())
				return "Failed to write data to storage!", err
			}
			logger.Info("Created new character for: %s", userID)
			return "Success!", nil
		}
	} else {
		logger.Error("Missing User ID from context!")
		return "", errors.New("Missing User ID from context!")
	}
}
