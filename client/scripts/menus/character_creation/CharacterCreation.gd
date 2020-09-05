extends Node

export(NodePath) var nameTextEdit
export(NodePath) var factionButtonsNode

onready var nameEdit : LineEdit = get_node(nameTextEdit)
onready var factionButtons = get_node(factionButtonsNode)

func _on_Button_button_down():
	if nameEdit.text != "":
		if factionButtons.currently_selected:
			for i in range(1,len(GameData.factions)):
				if GameData.factions[i] == factionButtons.currently_selected:
					var created_char = yield(ServerConnection.create_character_async(nameEdit.text, i), "completed")
					if created_char:
						get_tree().change_scene("res://scenes/World.tscn")
