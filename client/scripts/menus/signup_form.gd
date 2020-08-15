extends Popup

export(NodePath) var usernamePath
export(NodePath) var passwordPath
export(NodePath) var confirmPasswordPath
export(NodePath) var buttonPath
export(NodePath) var errorPath

var usernameEdit : LineEdit
var passwordEdit : LineEdit
var cPasswordEdit : LineEdit
var errorLabel : Label
var button : Button

const MIN_PASSWORD_LENGTH = 8

func _ready():
	# Get nodes
	usernameEdit = get_node(usernamePath)
	passwordEdit = get_node(passwordPath)
	cPasswordEdit = get_node(confirmPasswordPath)
	errorLabel = get_node(errorPath)
	button = get_node(buttonPath)
	
	# Set forms to validate on value chagne
	usernameEdit.connect("text_changed", self, "validate_fields")
	passwordEdit.connect("text_changed", self, "validate_fields")
	cPasswordEdit.connect("text_changed", self, "validate_fields")
	
	usernameEdit.connect("text_entered", self, "signup")
	passwordEdit.connect("text_entered", self, "signup")
	cPasswordEdit.connect("text_entered", self, "signup")
	
	# Connect submission button
	button.connect("button_down", self, "signup")
	
	# Clear error message
	errorLabel.text = ""

func signup(_text=""):
	if button.disabled:
		return
	
	var error : NakamaException = yield(ServerConnection.signup_async(usernameEdit.text, passwordEdit.text), "completed")
	
	# Check for error
	if error:
		passwordEdit.text = ""
		cPasswordEdit.text = ""
		errorLabel.add_color_override("font_color", Color.red)
		errorLabel.text = error.message
	else:
		errorLabel.add_color_override("font_color", Color.green)
		errorLabel.text = "Signed up successfully!"
		print("Signed up successfully!")
		hide()

func validate_fields(_text=""):
	var valid : bool = check_email(usernameEdit.text) and passwords_valid(passwordEdit.text, cPasswordEdit.text)
	button.disabled = !valid
	return valid

func passwords_valid(password, cpassword):
	return password == cpassword and len(password) >= MIN_PASSWORD_LENGTH

func check_email(email) -> bool:
	# Use regex to validate email
	var regex = RegEx.new()
	regex.compile("[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,6}")

	var result = regex.search(email)

	if result:
		return true
	return false
