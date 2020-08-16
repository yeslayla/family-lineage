extends Node

export(NodePath) var usernamePath
export(NodePath) var passwordPath
export(NodePath) var buttonPath
export(NodePath) var errorPath

var usernameEdit : LineEdit
var passwordEdit : LineEdit
var errorLabel : Label
var button : Button

func _ready():
	# Get nodes
	usernameEdit = get_node(usernamePath)
	passwordEdit = get_node(passwordPath)
	errorLabel = get_node(errorPath)
	button = get_node(buttonPath)
	
	# Connect submission button
	button.connect("button_down", self, "login")
	usernameEdit.connect("text_entered", self, "login")
	passwordEdit.connect("text_entered", self, "login")
	
	# Clear error message
	errorLabel.text = ""

func login(_text=""):
	var error : NakamaException = yield(ServerConnection.authenticate_async(usernameEdit.text, passwordEdit.text), "completed")
	
	# Check for error
	if error:
		passwordEdit.text = ""
		display_message(error.message)
	else:
		display_message("Logged in successfully!", Color.green)
		display_message("Connecting to server...", Color.gray)
		error = yield(ServerConnection.connect_to_server_async(), "completed")
		if error:
			display_message(error.message)
		else:
			display_message("Connected to server!", Color.green)
			# Load World

func display_message(message="", color=Color.red):
	errorLabel.add_color_override("font_color", color)
	errorLabel.text = message
	print(message)
