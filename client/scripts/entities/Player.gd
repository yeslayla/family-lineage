extends KinematicBody2D

signal finished_moving

export var base_movement_speed = 100
export(NodePath) var world

onready var navigation : Navigation2D = get_node(world)
var path : PoolVector2Array
var move : bool = false

func _process(delta):
	
	# Click for movement target
	if Input.is_action_just_pressed("move_to_cursor"):
		path = navigation.get_simple_path(global_position, get_global_mouse_position())
		move = true
	
	if move:
		move_along_path(base_movement_speed * delta)

func move_along_path(distance : float):
	var start_point := global_position
	for i in range(path.size()):
		var distance_to_next := start_point.distance_to(path[0])
		if distance <= distance_to_next and distance >= 0.0:
			global_position = start_point.linear_interpolate(path[0], distance/distance_to_next)
			break
		elif distance < 0.0:
			global_position = path[0]
			emit_signal("finished_moving")
			print("DONE")
			move = false
			break
		distance -= distance_to_next
		start_point = path[0]
		path.remove(0)

