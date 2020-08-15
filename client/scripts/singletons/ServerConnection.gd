extends Node

const KEY := "defaultkey"
const SERVER_ENDPOINT := "nakama.cloudsumu.com"

var _session : NakamaSession
var _client : NakamaClient = Nakama.create_client(KEY, SERVER_ENDPOINT, 7350, "http")

func authenticate_async(email : String, password : String) -> NakamaException:
	var result : NakamaException = null
	
	var new_session : NakamaSession = yield(_client.authenticate_email_async(email, password, null, false), "completed")
	
	if not new_session.is_exception():
		_session = new_session
	else:
		result = new_session.get_exception()
		
	return result
	
func signup_async(email : String, password : String) -> NakamaException:
	var result : NakamaException = null
	
	var new_session : NakamaSession = yield(_client.authenticate_email_async(email, password, null, true), "completed")
	
	if not new_session.is_exception():
		_session = new_session
	else:
		result = new_session.get_exception()
		
	return result
