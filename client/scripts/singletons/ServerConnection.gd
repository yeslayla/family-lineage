extends Node

signal tile_update(tile_data)

const KEY := "defaultkey"
const SERVER_ENDPOINT := "nakama.cloudsumu.com"

var _session : NakamaSession
var _client : NakamaClient = Nakama.create_client(KEY, SERVER_ENDPOINT, 7350, "http")
var _socket : NakamaSocket
var _precenses : Dictionary = {}

enum OPCODE {
	tile_update = 1
}

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
		_socket.connect("received_match_state", self, "_on_socket_received_match_state")
		_socket.connect("closed", self, "_on_socket_closed")
		return null
	return result.exception
	
func join_world_async() -> Dictionary:
	var world : NakamaAPI.ApiRpc = yield(_client.rpc_async(_session, "get_world_id", ""), "completed")
	if world.is_exception():
		print("Join world error occured: %s" % world.exception.message)
		return {}
		
	var match_join_result : NakamaRTAPI.Match = yield(_socket.join_match_async(world.payload), "completed")
	if match_join_result.is_exception():
		print("Join match error: %s - %s" % [match_join_result.exception.status_code, match_join_result.exception.message])
		return {}
	
	for precense in match_join_result.presences:
		_precenses[precense.user_id] = precense
		
	print("Currently connected: %s" % _precenses.size())

	return _precenses

func _on_socket_closed():
	_socket = null
	
func _on_socket_received_match_state(match_state: NakamaRTAPI.MatchData):
	match match_state.op_code:
		OPCODE.tile_update:
			emit_signal("tile_update", JSON.parse(match_state.data).result)

