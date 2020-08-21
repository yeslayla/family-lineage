package gamemap

import (
	"encoding/json"
	"fmt"
)

// WorldMap is the data structure used game world
type WorldMap struct {
	data  [256][256]int
	max_x int
	max_y int
}

// GetTile method is used to grab a tile value with error checking
func (m WorldMap) GetTile(x int, y int) (int, error) {
	if x > m.max_x || y > m.max_y {
		return -1, fmt.Errorf("Map out of bounds error: %d, %d", x, y)
	}
	return m.data[x][y], nil
}

// GetJSONRegion method returns a JSON object containing the tile values of everything
// within a given range
func (m WorldMap) GetJSONRegion(startX, endX, startY, endY int) ([]byte, error) {
	regionMap := map[int]map[int]int{}
	for x := startX; x < endX; x++ {
		regionMap[x] = map[int]int{}
		for y := startY; y < endY; y++ {

			// GetTile and ignore out of bounds errors
			result, _ := m.GetTile(x, y)
			regionMap[x][y] = result
		}
	}

	jsonString, err := json.Marshal(regionMap)
	return jsonString, err
}

// GetJSONRegionAround returns a JSON object of tile data from a center point
func (m WorldMap) GetJSONRegionAround(centerX float64, centerY float64, regionRadius int) ([]byte, error) {
	var xCenter int = int(centerX)
	var yCenter int = int(centerY)
	jsonString, err := m.GetJSONRegion(xCenter-regionRadius, xCenter+regionRadius, yCenter-regionRadius, yCenter+regionRadius)
	return jsonString, err
}

// IntializeMap is a method that helps easily
// generate WorldMap objects
func IntializeMap() *WorldMap {
	worldMap := new(WorldMap)
	worldMap.max_x = 256
	worldMap.max_y = 256
	worldMap.data = [256][256]int{}
	for x := 0; x < worldMap.max_x; x++ {
		for y := 0; y < worldMap.max_y; y++ {
			worldMap.data[x][y] = 0
		}
	}

	// Add dot at top left for testing
	worldMap.data[0][0] = 1

	return worldMap
}
