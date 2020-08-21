extends Node2D

export(NodePath) var puppet_parent
export(Resource) var puppet_template

onready var puppet_parent_node : Node = get_node(puppet_parent)

var puppets : Dictionary = {}

func _ready():
	ServerConnection.connect("player_joined", self, "on_player_join")
	ServerConnection.connect("player_left", self, "on_player_leave")
	ServerConnection.connect("player_pos_update", self, "on_player_pos_update")


func on_player_join(user_id):
	if user_id != ServerConnection._session.user_id:
		var new_puppet : Node = puppet_template.instance()
		new_puppet.name = "Player: " + user_id
		puppet_parent_node.add_child(new_puppet)
		puppets[user_id] = new_puppet

func on_player_leave(user_id):
	if user_id != ServerConnection._session.user_id:
		var player_puppet : Node = puppets[user_id]
		player_puppet.queue_free()

func on_player_pos_update(user_id, pos):
	if user_id != ServerConnection._session.user_id:
		var player_puppet : Node2D = puppets[user_id]
		player_puppet.global_position = pos
