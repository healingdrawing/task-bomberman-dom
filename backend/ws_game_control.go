package main

import (
	"log"
	"time"
)

func ws_character_control_handler(client *Client, control string) {

	switch control {
	case string(WS_UP_ON):
		ws_up_handler(string_number[client.NUMBER], control, true)
	case string(WS_UP_OFF):
		ws_up_off_handler(string_number[client.NUMBER], control)
	case string(WS_DOWN_ON):
		ws_down_handler(string_number[client.NUMBER], control, true)
	case string(WS_DOWN_OFF):
		ws_down_off_handler(string_number[client.NUMBER], control)
	case string(WS_LEFT_ON):
		ws_left_handler(string_number[client.NUMBER], control, true)
	case string(WS_LEFT_OFF):
		ws_left_off_handler(string_number[client.NUMBER], control)
	case string(WS_RIGHT_ON):
		ws_right_handler(string_number[client.NUMBER], control, true)
	case string(WS_RIGHT_OFF):
		ws_right_off_handler(string_number[client.NUMBER], control)
	case string(WS_BOMB_ON):
		go ws_bomb_handler(string_number[client.NUMBER]) // goroutine because plan to use time.Sleep, before bomb will be exploded
	case string(WS_BOMB_OFF):
		ws_bomb_off_handler(string_number[client.NUMBER])
	default:
		log.Println("Unknown control: ", control)
	}
}

// move_cells_per_second, depends on turbo
var move_cps = map[bool]int64{
	true:  2, // turbo
	false: 1,
}

func unpress_arrow(player *PLAYER, control string) {
	switch control {
	case string(WS_UP_OFF):
		player.up_pressed = false
	case string(WS_DOWN_OFF):
		player.down_pressed = false
	case string(WS_LEFT_OFF):
		player.left_pressed = false
	case string(WS_RIGHT_OFF):
		player.right_pressed = false
	default:
		log.Println("Unknown control to unpress: ", control)
	}
}

func unpress_all_arrows(player *PLAYER) {
	player.up_pressed = false
	player.down_pressed = false
	player.left_pressed = false
	player.right_pressed = false
}

// todo: check this . pointers or without.
func press_arrow_unpress_other_arrows(player *PLAYER, control string) {
	switch control {
	case string(WS_UP_ON):
		player.down_pressed = false
		player.left_pressed = false
		player.right_pressed = false
		player.up_pressed = true
	case string(WS_DOWN_ON):
		player.up_pressed = false
		player.left_pressed = false
		player.right_pressed = false
		player.down_pressed = true
	case string(WS_LEFT_ON):
		player.up_pressed = false
		player.down_pressed = false
		player.right_pressed = false
		player.left_pressed = true
	case string(WS_RIGHT_ON):
		player.up_pressed = false
		player.down_pressed = false
		player.left_pressed = false
		player.right_pressed = true
	default:
		log.Println("Unknown control to unpress arrows: ", control)
	}
}

// executed as goroutine. listen for player press arrows
func ws_arrows_loop_listener() {
	for game_waiting_state == GAME_STARTED {
		//sleep 0.1 sec
		time.Sleep(100 * time.Millisecond)
		game.Players.Range(func(key, value interface{}) bool {
			number := key.(string)
			player := value.(PLAYER)
			if !player.Dead {
				ws_up_handler(number, string(WS_UP_ON), false)
				ws_down_handler(number, string(WS_DOWN_ON), false)
				ws_left_handler(number, string(WS_LEFT_ON), false)
				ws_right_handler(number, string(WS_RIGHT_ON), false)
			}
			return true
		})
	}
}
