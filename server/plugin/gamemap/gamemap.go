package gamemap

import (
	"encoding/json"
	"fmt"
)

type WorldMap struct {
	data  [][]int
	max_x int
	max_y int
}

func (m WorldMap) GetTile(x int, y int) (int, error) {
	if x > m.max_x || y > m.max_y {
		return -1, fmt.Errorf("Map out of bounds error: %d, %d", x, y)
	}
	return m.data[x][y], nil
}

func (m WorldMap) GetJsonRegion(start_x, end_x, start_y, end_y int) ([]byte, error) {
	regionMap := make(map[int]map[int]int)
	for x := start_x; x < end_x; x++ {

		for y := start_y; y < end_y; y++ {
			if result, err := m.GetTile(x, y); err != nil {
				return nil, err
			} else {
				regionMap[x][y] = result
			}
		}
	}

	jsonString, err := json.Marshal(regionMap)
	return jsonString, err
}

func IntializeMap() *WorldMap {
	worldMap := new(WorldMap)
	worldMap.max_x = 64
	worldMap.max_y = 64
	for x := 0; x < worldMap.max_x; x++ {
		for y := 0; y < worldMap.max_y; y++ {
			worldMap.data[x][y] = 0
		}
	}

	// Add dot at top left for testing
	worldMap.data[0][0] = 1

	return worldMap
}
