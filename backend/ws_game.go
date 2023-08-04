package main

type WS_START_GAME_DTO struct {
	Content string `json:"content"`
}

func ws_send_start_game_handler() {
	uuids := get_all_clients_uuids(clients)
	wsSend(WS_START_GAME, WS_START_GAME_DTO{"gap. Must be full initial state for game"}, uuids)
}
