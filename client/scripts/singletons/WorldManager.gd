extends Node

export(NodePath) var tilemapPath

var tilemap : TileMap

func _ready():
	
	# Setup tilemap
	tilemap = get_node(tilemapPath)
	tilemap.clear()
	
	# Setup connections and join wolrd
	ServerConnection.connect("tile_update", self, "on_tile_update")
	yield(ServerConnection.join_world_async(), "completed")

func on_tile_update(tile_data, update_bitmask=true):
	print("Updating tilemap")
	
	var max_pos_x : int
	var min_pos_x : int
	var max_pos_y : int
	var min_pos_y : int
	
	for x in tile_data:
		
		# Find max & min x
		if not max_pos_x or max_pos_x > int(x):
			max_pos_x = int(x)
		if not min_pos_x or min_pos_x < int(x):
			min_pos_x = int(x)
		
		for y in tile_data[x]:
			
			# Find max & min y
			if not max_pos_y or max_pos_y > int(y):
				max_pos_y = int(y)
			if not min_pos_y or min_pos_y < int(y):
				min_pos_y = int(y)
			
			# Update tile data
			tilemap.set_cell(int(x),int(y), int(tile_data[x][y]), false, false, false, tilemap.get_cell_autotile_coord(int(x), int(y)))
	
	if update_bitmask:
		tilemap.update_bitmask_region(Vector2(min_pos_x, min_pos_y), Vector2(max_pos_x, max_pos_y))

	print("Update complete!")
