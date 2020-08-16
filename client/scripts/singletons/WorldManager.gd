extends Node


# Declare member variables here. Examples:
# var a = 2
# var b = "text"


# Called when the node enters the scene tree for the first time.
func _ready():
	ServerConnection.connect("tile_update", self, "on_tile_update")
	yield(ServerConnection.join_world_async(), "completed")

func on_tile_update(tile_data):
	print(tile_data)
