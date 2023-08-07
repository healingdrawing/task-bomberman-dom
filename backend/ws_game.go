package main

import "sync"

var (
	game       = GAME_STATE{}
	game_mutex = sync.Mutex{}
)

type PLAYER struct {
	Number                  int   `json:"number"` // the Client.NUMBER
	X                       int   `json:"x"`      // the number of the game field cell, where the player is
	Y                       int   `json:"y"`
	Target_x                int   `json:"target_x"` // the number of the game field cell, where the player moves to
	Target_y                int   `json:"target_y"`
	bombs_max               int   // default 1, can be increased by powerup. how many bombs can be placed at the same time
	bombs_left              int   // how many bombs can be placed at the moment
	explosion_range         int   // default 1, can be increased by powerup
	Turbo                   bool  `json:"turbo"` // default false, can be switched by powerup. speedup of player
	Dead                    bool  `json:"dead"`  // true - player is dead, false - player is alive
	moving                  bool  // true - player is moving, false - player is not moving
	moving_start_time_stamp int64 // unix timestamp in nanoseconds, when the player started moving
	one_cell_move_duration  int64 // how many nanoseconds does it take to move one cell
	up_pressed              bool
	down_pressed            bool
	left_pressed            bool
	right_pressed           bool
	bomb_pressed            bool
}

type POWER_UP struct {
	Effect string `json:"effect"` // "bombs_max", "explosion_range", "turbo"
	Show   bool   `json:"show"`   // true - show on the game field, false - already taken by player or under the weak obstacle
}

type GAME_STATE struct {
	Players        map[string]PLAYER   `json:"players"`
	Weak_obstacles map[string]bool     `json:"weak_obstacles"` // key is xy, without space, like "01", value is show(true) or destroyed(false)
	Power_ups      map[string]POWER_UP `json:"power_ups"`      // key is xy, without space, like "01"
	free_cells     map[string]bool     // game field , except weak and strong obstacles cellxy. can move there. key is xy, without space, like "01"
}

func game_init() {

	prepare_players()

	prepare_weak_obstacles_and_power_ups()

	prepare_free_cells() // strong and weak obstacles to restrict movement

	go ws_arrows_loop_listener()

}

func ws_send_start_game_handler() {
	game_init()

	uuids := get_all_clients_uuids(clients)
	wsSend(WS_START_GAME, game, uuids)
}
