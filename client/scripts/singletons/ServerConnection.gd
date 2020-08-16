extends Node

const KEY := "defaultkey"
const SERVER_ENDPOINT := "nakama.cloudsumu.com"

var _session : NakamaSession
var _client : NakamaClient = Nakama.create_client(KEY, SERVER_ENDPOINT, 7350, "http")
var _socket : NakamaSocket

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
	
	var new_session : NakamaSession = yield(_client.authenticate_email_async(email, password, email, true), "completed")
	
	if not new_session.is_exception():
		_session = new_session
	else:
		result = new_session.get_exception()
		
	return result

func connect_to_server_async() -> NakamaException:
	_socket = Nakama.create_socket_from(_client)
	var result : NakamaAsyncResult = yield(_socket.connect_async(_session), "completed")
	if not result.is_exception():
		_socket.connect("closed", self, "_on_socket_closed")
		return null
	return result.exception
	
func join_world_async() -> Dictionary:
	var world : NakamaAPI.ApiRpc = yield(_client.rpc_async(_session, "get_world_id", ""), "completed")
	if world.is_exception():
		print("Join world error occured: %s" % world.exception.message)
		return {}
	var _world_id : String = world.payload
	print(_world_id)
	return {}

func _on_socket_closed():
	_socket = null
