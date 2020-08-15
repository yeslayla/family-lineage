extends Node


# Declare member variables here. Examples:
# var a = 2
# var b = "text"


# Called when the node enters the scene tree for the first time.
func _ready():
	return
	var result : int
	result = yield(ServerConnection.authenticate_async("j@cloudsumu.com", "Learning12!"), "completed")
	if result == null:
		print("Logged In")
	else:
		print("Auth failed! Error code: %d" % result)
		result = yield(ServerConnection.signup_async("j@cloudsumu.com", "Learing12!"), "completed")
		if result == null:
			print("Registered!")
		else:
			print("Signup failed! Error code: %d" % result)
