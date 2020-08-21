package entities

import (
	"encoding/json"
	"fmt"
	"github.com/heroiclabs/nakama-common/runtime"
	"strconv"
)

// PlayerEntity is the go struct representing the player's location
type PlayerEntity struct {
	Presence runtime.Presence
	X        float64
	Y        float64
}

// PlayerPosResponse struct that represents client data
type PlayerPosResponse struct {
	X string
	Y string
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
