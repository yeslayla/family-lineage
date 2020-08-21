extends Node

signal tile_update(tile_data)
signal player_joined(user_id)
signal player_left(user_id)
signal player_pos_update(user_id, pos)

const KEY := "defaultkey"
const SERVER_ENDPOINT := "nakama.cloudsumu.com"

var _session : NakamaSession
var _client : NakamaClient = Nakama.create_client(KEY, SERVER_ENDPOINT, 7350, "http")
var _socket : NakamaSocket
var _precenses : Dictionary = {}
var _world_id

enum OPCODE {
	tile_update = 1,
	update_position = 2
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
		_socket.connect("received_match_presence", self, "_on_received_match_presence")
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
		
	_world_id = world.payload
	
	for precense in match_join_result.presences:
		_precenses[precense.user_id] = precense
		emit_signal("player_joined", precense.user_id)
		
	print("Joined matched with %s other players!" % _precenses.size())

	return _precenses

func _on_socket_closed():
	_socket = null

func _on_received_match_presence(match_precense : NakamaRTAPI.MatchPresenceEvent):
	for precense in match_precense.joins:
		print("%s joined the game!" % precense.username)
		_precenses[precense.user_id] = precense
		emit_signal("player_joined", precense.user_id)
		
	for precense in match_precense.leaves:
		print("%s left the game!" % precense.username)
		_precenses.erase(precense.user_id)
		emit_signal("player_left", precense.user_id)

func _on_socket_received_match_state(match_state: NakamaRTAPI.MatchData):
	match match_state.op_code:
		OPCODE.tile_update:
			emit_signal("tile_update", JSON.parse(match_state.data).result)
		OPCODE.update_position:
			var pos_data = JSON.parse(match_state.data).result
			emit_signal("player_pos_update", pos_data["player"], Vector2(float(pos_data["x"]), float(pos_data["y"])))

func send_player_position(position : Vector2) -> void:
	_socket.send_match_state_async(_world_id, OPCODE.update_position, JSON.print({X = str(position.x), Y = str(position.y)}))

