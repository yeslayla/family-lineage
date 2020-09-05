extends Node

signal selection_updated(faction)

var currently_selected : String

func _ready():
	get_child(0).queue_free()
	for faction in GameData.factions:
		if faction == "Electus":
			return
		
		var faction_button = Button.new()
		faction_button.icon = load("res://art/gui/banners/" + faction.to_lower() + ".png")
		faction_button.connect("button_down", self, "on_faction_select", [faction])
		add_child(faction_button)

func on_faction_select(faction):
	currently_selected = faction
	emit_signal("selection_updated", faction)
