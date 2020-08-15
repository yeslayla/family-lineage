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
		errorLabel.add_color_override("font_color", Color.red)
		errorLabel.text = error.message
	else:
		errorLabel.add_color_override("font_color", Color.green)
		errorLabel.text = "Logged in successfully!"
		print("Logged in successfully!")
