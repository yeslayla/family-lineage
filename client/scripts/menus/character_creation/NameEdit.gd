extends LineEdit

var invalid_chars = " '.,;/\\,[](){}!@#$%^&*-=|_+1234567890\""

func _ready():
	connect("text_changed", self, "validate_name")


func validate_name(name):
	var cursor_pos = caret_position
	for character in invalid_chars:
		if character in text:
			text = text.replace(character, "")
			caret_position = cursor_pos
