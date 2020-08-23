package control

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/josephbmanley/family/server/plugin/entities"
	"github.com/josephbmanley/family/server/plugin/gamemap"
)

const maxRenderDistance int = 32

// OpCode represents a enum for valid OpCodes
// used by the match logic
type OpCode int64

const (
	// OpCodeTileUpdate is used for tile updates
	OpCodeTileUpdate = 1
	// OpCodeUpdatePosition is used for player position updates
	OpCodeUpdatePosition = 2
	// OpCodePlayerState is for player object updates
	OpCodePlayerState = 3
)

// Match is the object registered
// as a runtime.Match interface
type Match struct{}

// MatchState holds information that is passed between
// Nakama match methods
type MatchState struct {
	presences map[string]runtime.Presence
	players   map[string]entities.PlayerEntity
	inputs    map[string]string
	worldMap  *gamemap.WorldMap
}

// GetPrecenseList returns an array of current precenes in an array
func (state *MatchState) GetPrecenseList() []runtime.Presence {
	precenseList := []runtime.Presence{}
	for _, precense := range state.presences {
		precenseList = append(precenseList, precense)
	}
	return precenseList
}

// MatchInit is called when a new match is created
func (m *Match) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {

	state := &MatchState{
		presences: map[string]runtime.Presence{},
		inputs:    map[string]string{},
		players:   map[string]entities.PlayerEntity{},
		worldMap:  gamemap.IntializeMap(),
	}
	tickRate := 10
	label := "{\"name\": \"Game World\"}"

	return state, tickRate, label
}

// MatchJoinAttempt is called when a player tried to join a match
// and validates their attempt
func (m *Match) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	mState, ok := state.(*MatchState)
	if !ok {
		logger.Error("Invalid match state on join attempt!")
		return state, false, "Invalid match state!"
	}

	// Validate user is not already connected
	if _, ok := mState.presences[presence.GetUserId()]; ok {
		return mState, false, "User already logged in."
	} else {

		dataExist, err := entities.PlayerDataExists(ctx, nk, presence.GetUserId())
		if err != nil {
			logger.Error(err.Error())
			return mState, false, err.Error()
		}

		if dataExist {
			return mState, true, ""
		} else {
			return mState, false, "User does not have a character!"
		}

	}

}

// MatchJoin is called when a player successfully joins the match
func (m *Match) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	mState, ok := state.(*MatchState)
	if !ok {
		logger.Error("Invalid match state on join!")
		return state
	}

	for _, precense := range presences {

		// Add presence to map
		mState.presences[precense.GetUserId()] = precense

		player, loadPlayerErr := entities.LoadPlayer(ctx, nk, precense)
		if loadPlayerErr != nil {
			logger.Error(loadPlayerErr.Error())
			player.X = 16
			player.Y = 16
		} else {
			logger.Error(loadPlayerErr.Error())
			player = entities.PlayerEntity{
				X:        16,
				Y:        16,
				Name:     "ERROR",
				Presence: precense,
			}
		}

		if jsonObj, err := player.GetPosJSON(); err != nil {
			logger.Error(err.Error())
		} else {
			if sendErr := dispatcher.BroadcastMessage(OpCodeUpdatePosition, jsonObj, []runtime.Presence{precense}, player.Presence, true); sendErr != nil {
				logger.Error(sendErr.Error())
			}
			stateJSON, _ := player.GetStateJSON()
			if sendErr := dispatcher.BroadcastMessage(OpCodePlayerState, stateJSON, mState.GetPrecenseList(), player.Presence, true); sendErr != nil {
				logger.Error(sendErr.Error())
			}
		}

		mState.players[precense.GetUserId()] = player

		// Get intial tile data around player
		if regionData, err := mState.worldMap.GetJSONRegionAround(player.X, player.Y, maxRenderDistance); err != nil {
			logger.Error(err.Error())
		} else {

			// Broadcast tile data to client
			if sendErr := dispatcher.BroadcastMessage(OpCodeTileUpdate, regionData, []runtime.Presence{precense}, precense, true); sendErr != nil {
				logger.Error(sendErr.Error())
			}
		}
		for _, otherPlayer := range mState.players {
			// Broadcast player data to client
			if jsonObj, err := otherPlayer.GetPosJSON(); err != nil {
				logger.Error(err.Error())
			} else {
				if sendErr := dispatcher.BroadcastMessage(OpCodeUpdatePosition, jsonObj, []runtime.Presence{precense}, otherPlayer.Presence, true); sendErr != nil {
					logger.Error(sendErr.Error())
				}
			}
		}
	}
	return mState
}

// MatchLeave is called when a player leaves the match
func (m *Match) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	mState, ok := state.(*MatchState)
	if !ok {
		logger.Error("Invalid match state on leave!")
		return state
	}
	for _, presence := range presences {
		if _, ok := mState.presences[presence.GetUserId()]; ok {
			delete(mState.presences, presence.GetUserId())
		}
		if _, ok := mState.players[presence.GetUserId()]; ok {
			delete(mState.players, presence.GetUserId())
		}
	}
	return mState
}

// MatchLoop is code that is executed every tick
func (m *Match) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	mState, ok := state.(*MatchState)
	if !ok {
		logger.Error("Invalid match state on leave!")
		return state
	}
	for _, message := range messages {
		if message.GetOpCode() == OpCodeUpdatePosition {
			player := mState.players[message.GetUserId()]

			if response, err := player.ParsePositionRequest(message.GetData()); err == nil {
				player.UpdateBasedOnResponse(response)
				if jsonObject, err := player.GetPosJSON(); err == nil {
					dispatcher.BroadcastMessage(OpCodeUpdatePosition, jsonObject, mState.GetPrecenseList(), player.Presence, false)
				} else {
					logger.Error(fmt.Sprintf("Failed to get player json: %s", err.Error))
				}
			} else {
				logger.Error(fmt.Sprintf("Failed to parse update pos request: %s", err.Error))
			}
		}
	}
	return mState
}

// MatchTerminate is code that is executed when the match ends
func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	return state
}
