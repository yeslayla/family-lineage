package entities

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/josephbmanley/family/server/plugin/gameworld"
	"strconv"
)

type PlayerSaveData struct {
	Faction int
	Name    string
}

// PlayerEntity is the go struct representing the player's location
type PlayerEntity struct {
	Presence runtime.Presence
	Faction  gameworld.Faction
	Name     string
	X        float64
	Y        float64
}

// PlayerPosResponse struct that represents client data
type PlayerPosResponse struct {
	X string
	Y string
}

// PlayerDataExists checks if precense has saved data
func PlayerDataExists(ctx context.Context, nk runtime.NakamaModule, presence runtime.Presence) (bool, error) {

	Reads := []*runtime.StorageRead{
		&runtime.StorageRead{
			Collection: "playerdata",
			Key:        "data",
			UserID:     presence.GetUserId(),
		},
	}
	records, err := nk.StorageRead(ctx, Reads)
	if err != nil {
		return false, err
	}
	for _, record := range records {
		if record.Key == "data" {
			return true, nil
		}
	}

	return false, nil
}

// LoadPlayer creates player object
func LoadPlayer(ctx context.Context, nk runtime.NakamaModule, presence runtime.Presence) (PlayerEntity, error) {
	player := PlayerEntity{Presence: presence}

	// Read storage
	PlayerReads := []*runtime.StorageRead{
		&runtime.StorageRead{
			Collection: "playerdata",
			Key:        "data",
			UserID:     player.Presence.GetUserId(),
		},
	}
	records, err := nk.StorageRead(ctx, PlayerReads)
	if err != nil {
		return player, err
	}

	// Load storage records into object
	for _, record := range records {
		switch record.Key {
		case "data":
			responseData := PlayerSaveData{}
			err := json.Unmarshal([]byte(record.Value), &responseData)
			if err != nil {
				return player, err
			}
			player.Name = responseData.Name
			player.Faction = gameworld.Faction(responseData.Faction)
		}
	}
	return player, nil
}

// Save saves player data to nakama
func (p *PlayerEntity) Save(ctx context.Context, nk runtime.NakamaModule) error {

	saveData := PlayerSaveData{
		Name:    p.Name,
		Faction: int(p.Faction),
	}

	saveJSON, err := json.Marshal(saveData)
	if err != nil {
		return err
	}

	PlayerWrites := []*runtime.StorageWrite{
		&runtime.StorageWrite{
			Collection: "playerdata",
			Key:        "data",
			Value:      string(saveJSON),
			UserID:     p.Presence.GetUserId(),
		},
	}

	_, err = nk.StorageWrite(ctx, PlayerWrites)

	return nil
}

// ParsePositionRequest parses data from client
func (p *PlayerEntity) ParsePositionRequest(data []byte) (PlayerPosResponse, error) {
	var response PlayerPosResponse
	err := json.Unmarshal(data, &response)
	return response, err
}

//UpdateBasedOnResponse updates the player object based on a response object
func (p *PlayerEntity) UpdateBasedOnResponse(response PlayerPosResponse) error {
	if fx, err := strconv.ParseFloat(response.X, 64); err != nil {
		return err
	} else {
		p.X = fx
		if fy, err := strconv.ParseFloat(response.Y, 64); err != nil {
			return err
		} else {
			p.Y = fy
		}
	}

	return nil
}

// GetPosJSON returns the player's position as a JSON object
func (p *PlayerEntity) GetPosJSON() ([]byte, error) {
	playerMap := map[string]string{
		"player": p.Presence.GetUserId(),
		"x":      fmt.Sprintf("%f", p.X),
		"y":      fmt.Sprintf("%f", p.Y),
	}
	jsonData, err := json.Marshal(playerMap)
	return jsonData, err
}
