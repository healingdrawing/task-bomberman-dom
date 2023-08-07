package main

import (
	"fmt"
	"log"
	"time"
)

func ws_character_control_handler(client *Client, control string) {

	switch control {
	case string(WS_UP_ON):
		ws_up_handler(string_number[client.NUMBER], control)
	case string(WS_UP_OFF):
		log.Println("WS_UP_OFF", client.NUMBER)
	case string(WS_DOWN_ON):
		log.Println("WS_DOWN_ON", client.NUMBER)
	case string(WS_DOWN_OFF):
		log.Println("WS_DOWN_OFF", client.NUMBER)
	case string(WS_LEFT_ON):
		log.Println("WS_LEFT_ON", client.NUMBER)
	case string(WS_LEFT_OFF):
		log.Println("WS_LEFT_OFF", client.NUMBER)
	case string(WS_RIGHT_ON):
		log.Println("WS_RIGHT_ON", client.NUMBER)
	case string(WS_RIGHT_OFF):
		log.Println("WS_RIGHT_OFF", client.NUMBER)
	case string(WS_BOMB_ON):
		log.Println("WS_BOMB_ON", client.NUMBER)
	case string(WS_BOMB_OFF):
		log.Println("WS_BOMB_OFF", client.NUMBER)
	default:
		log.Println("Unknown control: ", control)
	}
}

// move_cells_per_second, depends on turbo
var move_cps = map[bool]int64{
	true:  2, // turbo
	false: 1,
}

func unpress_all_arrows(player *PLAYER) {
	player.up_pressed = false
	player.down_pressed = false
	player.left_pressed = false
	player.right_pressed = false
}

// todo: check this . pointers or without.
func unpress_other_arrows(player *PLAYER, control string) {
	switch control {
	case string(WS_UP_ON):
		player.down_pressed = false
		player.left_pressed = false
		player.right_pressed = false
	case string(WS_DOWN_ON):
		player.up_pressed = false
		player.left_pressed = false
		player.right_pressed = false
	case string(WS_LEFT_ON):
		player.up_pressed = false
		player.down_pressed = false
		player.right_pressed = false
	case string(WS_RIGHT_ON):
		player.up_pressed = false
		player.down_pressed = false
		player.left_pressed = false
	default:
		log.Println("Unknown control to unpress arrows: ", control)
	}
}

// executed as goroutine. listen for player press arrows
func ws_arrows_loop_listener() {
	for game_waiting_state == GAME_STARTED {
		game.Players.Range(func(key, value interface{}) bool {
			number := key.(string)
			player := value.(PLAYER)
			if !player.Dead && player.up_pressed {
				ws_up_handler(number, string(WS_UP_ON))
			}
			return true
		})
	}
}

// player press up arrow
func ws_up_handler(number string, control string) {
	log.Println("ws_up_handler. player", number)

	log.Println("step 0")

	playerValue, ok := game.Players.Load(number)
	if !ok {
		return
	}

	player := playerValue.(PLAYER)
	log.Println("step 1")

	unpress_other_arrows(&player, control)
	game.Players.Store(number, player)

	log.Println("step 2")
	log.Println("player.moving ", player.moving)

	// Unix timestamp in nanoseconds
	unix_ts := time.Now().UnixNano()

	if !player.moving {
		log.Println("step 2.1")
		// Check if player can move up
		target_position_cell := fmt.Sprintf("%d%d", player.X, player.Y-1)
		log.Println("step 3")
		if _, ok := game.free_cells.Load(target_position_cell); !ok {
			player.up_pressed = false
			return
		}
		log.Println("step 4")
		seconds := 1
		oneSecond := time.Duration(seconds) * time.Second
		// Time in nanoseconds
		one_cell_move_duration := oneSecond.Nanoseconds() / move_cps[player.Turbo]
		player.one_cell_move_duration = one_cell_move_duration
		player.moving_start_time_stamp = unix_ts
		player.Target_y = player.Y - 1
		player.moving = true
		log.Println("step 5")
		game.Players.Store(number, player)
		log.Println("step 6")
		ws_send_move_up_command(&player)
		log.Println("step 7")
	} else if player.one_cell_move_duration/10*5 < unix_ts-player.moving_start_time_stamp {
		// Check if 50% of the timer is done, switch to the next cell
		player.Y = player.Target_y
		game.Players.Store(number, player)
	} else if player.one_cell_move_duration/10*8 < unix_ts-player.moving_start_time_stamp {
		// Check if 80% of the timer is done, and the player still wants to move up
		// If yes, then move/aim the player up one cell more if possible
		// 100000 is 100 microseconds gap to wait while after timesleep is over, recursive call is done
		target_position_cell := fmt.Sprintf("%d%d", player.X, player.Y-1)
		if _, ok := game.free_cells.Load(target_position_cell); !ok {
			player.up_pressed = false
			return
		}
		player.Target_y = player.Y - 1
		player.moving_start_time_stamp += player.one_cell_move_duration
		game.Players.Store(number, player)
		ws_send_move_up_command(&player)
	}
}

func ws_send_move_up_command(player *PLAYER) {
	log.Println("ws_send_move_up_command. player", player.Number)
	message := WS_MOVE_UP_DTO{
		Number:   player.Number,
		Target_y: player.Target_y,
		Turbo:    player.Turbo,
	}
	log.Println("ws_send_move_up_command. message", message)
	uuids := get_all_clients_uuids(clients)
	log.Println("ws_send_move_up_command. uuids", uuids)
	wsSend(WS_UP, message, uuids)
	log.Println("ws_send_move_up_command. sent")

}

type WS_MOVE_UP_DTO struct {
	Number   int  `json:"number"`
	Target_y int  `json:"target_y"`
	Turbo    bool `json:"turbo"`
}
