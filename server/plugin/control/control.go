package control

import (
	"context"
	"database/sql"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/josephbmanley/family/server/plugin/gamemap"
)

type OpCode int64

const (
	OpCodeTileUpdate = 1
)

type Match struct{}

type MatchState struct {
	presences map[string]runtime.Presence
	inputs    map[string]string
	positions map[string]map[string]int
	names     map[string]string
	worldMap  *gamemap.WorldMap
}

func (m *Match) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {

	state := &MatchState{
		presences: map[string]runtime.Presence{},
		inputs:    map[string]string{},
		positions: map[string]map[string]int{},
		names:     map[string]string{},
		worldMap:  gamemap.IntializeMap(),
	}
	tickRate := 10
	label := "{\"name\": \"Game World\"}"

	return state, tickRate, label
}

func (m *Match) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	mState, ok := state.(*MatchState)
	if !ok {
		logger.Error("Invalid match state on join attempt!")
		return state, false, "Invalid match state!"
	}
	if _, ok := mState.presences[presence.GetUserId()]; ok {
		return mState, false, "User already logged in."
	} else {
		return mState, true, ""
	}

}

func (m *Match) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	mState, ok := state.(*MatchState)
	if !ok {
		logger.Error("Invalid match state on join!")
		return state, false, "Invalid match state!"
	}
	for _, precense := range presences {
		mState.presences[precense.GetUserId()] = precense

		mState.positions[precense.GetUserId()] = map[string]int{"x": 16, "y": 16}

		mState.names[precense.GetUserId()] = "User"

		if regionData, err := mState.worldMap.GetJsonRegion(16-8, 16+8, 16-8, 16+8); err != nil {
			logger.Error(err.Error())
			return mState
		} else {
			if sendErr := dispatcher.BroadcastMessage(OpCodeTileUpdate, regionData, []runtime.Presence{precense}, precense, true); sendErr != nil {
				logger.Error(sendErr.Error())
				return mState
			}
		}
	}
	return mState
}

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

func (m *Match) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	// Custom code to:
	// - Process the messages received.
	// - Update the match state based on the messages and time elapsed.
	// - Broadcast new data messages to match participants.

	return state
}

func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	return state
}
