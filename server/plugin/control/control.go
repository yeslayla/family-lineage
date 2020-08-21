package control

import (
	"context"
	"database/sql"
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
		return mState, true, ""
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

		player := entities.PlayerEntity{
			X: 16,
			Y: 16,
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
		delete(mState.presences, presence.GetUserId())
	}
	return mState
}

// MatchLoop is code that is executed every tick
func (m *Match) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	// Custom code to:
	// - Process the messages received.
	// - Update the match state based on the messages and time elapsed.
	// - Broadcast new data messages to match participants.

	return state
}

// MatchTerminate is code that is executed when the match ends
func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	return state
}
