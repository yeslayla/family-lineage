extends "res://addons/gut/test.gd"

var signup_form = load("res://scripts/menus/signup_form.gd")

# Test Object
var form = signup_form.new()

#------------
# Email Test
#------------
var valid_email_list = [
		"untitled@gmail.com",
		"test@cloudsumu.com",
		"cool.game@tetraforce.io",
		"ExampleName@yahoo.com"
]

var invalid_email_list = [
	"test the test",
	"test",
	"test@test",
	"gmail.com",
	"google.com",
	"@amazon.com",
	"test@_.com",
	"test@test.",
	"Hello World!"
]

func test_check_email_with_valid_email():
	for email in valid_email_list:
		assert_true(form.check_email(email))

func test_check_email_with_invalid_email():
	for email in invalid_email_list:
		assert_false(form.check_email(email))

#---------------
# Password Test
#---------------
var valid_passwords = [
	"Testing123!",
	"gR8$cuP8kJ8%qk*t",
	"GVa9%BZHh",
	"2Uw@2*5Qb$Gflb@c",
	"iL3DINd@hRaBlevo"
]

func test_passwords_valid_do_match():
	for password in valid_passwords:
		assert_true(form.passwords_valid(password, password))

func test_passwords_valid_do_not_match():
	for password in valid_passwords:
		assert_false(form.passwords_valid(password, null))
