package main

import (
	"log"
	"sync"
)

type DIRECTION int

const (
	STAND DIRECTION = 1 + iota
	UP
	DOWN
	LEFT
	RIGHT
)

var (
	game_debug = true
	game       = GAME_STATE{}
	// player numbers as string short hand
	string_number = []string{"0", "1", "2", "3", "4"}
	// number from string. at least to calculate explosion range, 0 to 6 inclusive
	cell_number_from_string = map[string]int{
		"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6,
	}
)

type PLAYER struct {
	Number                  int       `json:"number"` // the Client.NUMBER
	X                       int       `json:"x"`      // the number of the game field cell, where the player is
	Y                       int       `json:"y"`
	Target_x                int       `json:"target_x"` // the number of the game field cell, where the player moves to
	Target_y                int       `json:"target_y"`
	bombs_max               int       // default 1, can be increased by powerup. how many bombs can be placed at the same time
	bombs_used              int       // how many bombs placed at the moment on the game field
	explosion_range         int       // default 1, can be increased by powerup
	Turbo                   bool      `json:"turbo"` // default false, can be switched by powerup. speedup of player
	Dead                    bool      `json:"dead"`  // true - player is dead, false - player is alive
	direction               DIRECTION // the direction of the player movement
	moving                  bool      // true - player is moving, false - player is not moving
	moving_start_time_stamp int64     // unix timestamp in nanoseconds, when the player started moving
	one_cell_move_duration  int64     // how many nanoseconds does it take to move one cell
	up_pressed              bool
	down_pressed            bool
	left_pressed            bool
	right_pressed           bool
	bomb_pressed            bool

	can_change_present_cell bool // true - player present position was not changed
	can_change_target_cell  bool // true - player target position was not changed
}

type POWER_UP struct {
	Effect string `json:"effect"` // "bombs_max", "explosion_range", "turbo"
	Show   bool   `json:"show"`   // true - show on the game field, false - already taken by player or under the weak obstacle
}

type GAME_STATE struct {
	Players        sync.Map `json:"players"`
	Weak_obstacles sync.Map `json:"weak_obstacles"` // key is xy, without space, like "01", value is show(true) or destroyed(false)
	Power_ups      sync.Map `json:"power_ups"`      // key is xy, without space, like "01"
	free_cells     sync.Map // game field , except weak and strong obstacles cellxy. can move there. key is xy, without space, like "01"
	bomb_cells     sync.Map // key is xy, without space, like "01"
}

func game_init() {

	prepare_players()

	prepare_weak_obstacles_and_power_ups()

	prepare_free_cells_and_bomb_cells() // strong and weak obstacles to restrict movement

	go ws_arrows_loop_listener()

}

func ws_send_start_game_handler() {
	game_init()

	uuids := get_all_clients_uuids(clients)
	wsSend(WS_START_GAME, convert_game_state_to_json(game), uuids)
}

func dprint(msg ...any) {
	if game_debug {
		log.Println(msg...)
	}
}
