extends "res://addons/gut/test.gd"

var world_manager = load("res://scripts/singletons/WorldManager.gd")

func test_adding_tiles_to_map():
	# Configure world to have a 
	var world = world_manager.new()
	world.tilemap = TileMap.new()
	
	world.on_tile_update({
		"0" : {
			"0" : "0",
			"1" : "0",
			"2" : "0",
			"3" : "1",
			"4" : "0"
		},
		"1" : {
			"0" : "1"
		}
	}, false)
	
	assert_eq(world.tilemap.get_cell(0,0), 0)
	assert_eq(world.tilemap.get_cell(0,1), 0)
	assert_eq(world.tilemap.get_cell(0,2), 0)
	assert_eq(world.tilemap.get_cell(0,3), 1)
	assert_eq(world.tilemap.get_cell(0,4), 0)
	assert_eq(world.tilemap.get_cell(1,0), 1)
	
	# Test Updates
	world.on_tile_update({
		"0" : {
			"1" : "0",
			"2" : "1"
		}
	}, false)

	assert_eq(world.tilemap.get_cell(0,1), 0)
	assert_eq(world.tilemap.get_cell(0,2), 1)
	
	# Test New Additions
	world.on_tile_update({
		"1" : {
			"6" : "0",
			"7" : "1"
		}
	}, false)

	assert_eq(world.tilemap.get_cell(1,6), 0)
	assert_eq(world.tilemap.get_cell(1,7), 1)
