extends Sprite

export var faction = 1 setget set_faction, get_faction

func _ready():
	set_faction(faction)

func set_faction(new_faction):
	faction = new_faction
	texture = load("res://art/entities/player/dev/dev_player_" + str(faction) + ".png")
	
func get_faction():
	return faction
