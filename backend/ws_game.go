package main

var (
	game = GAME_STATE{}
)

type PLAYER struct {
	number          int // the Client.NUMBER
	x               int // the number of the game field cell, where the player is
	y               int
	target_x        int // the number of the game field cell, where the player moves to
	target_y        int
	bombs_max       int  // default 1, can be increased by powerup. how many bombs can be placed at the same time
	bombs_left      int  // how many bombs can be placed at the moment
	explosion_range int  // default 1, can be increased by powerup
	turbo           bool // default 1, can be increased by powerup. speedup of player
	dead            bool
}

type WEAK_OBSTACLE struct {
	x int
	y int
}

type GAME_STATE struct {
	players        map[string]PLAYER
	weak_obstacles map[string]WEAK_OBSTACLE // key is xy, without space, like "01"
}

func game_init() {

	prepare_players()

	prepare_weak_obstacles()

}

type WS_START_GAME_DTO struct {
	Content string `json:"content"`
}

func ws_send_start_game_handler() {
	game_init()

	uuids := get_all_clients_uuids(clients)
	wsSend(WS_START_GAME, WS_START_GAME_DTO{"gap. Must be full initial state for game"}, uuids)
}
