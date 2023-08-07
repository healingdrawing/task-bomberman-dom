package main

import (
	"fmt"
	"log"
	"time"
)

func ws_character_control_handler(client *Client, control string) {

	switch control {
	case string(WS_UP_ON):
		ws_up_on_handler(client.NUMBER, control)
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

// player numbers short hand
var pn = []string{"0", "1", "2", "3", "4"}

// move_cells_per_second, depends on turbo
var move_cps = map[bool]int64{
	true:  2, // turbo
	false: 1,
}

func unpress_all_arrows(player PLAYER) {
	player.up_pressed = false
	player.down_pressed = false
	player.left_pressed = false
	player.right_pressed = false
}

func unpress_other_arrows(player PLAYER, control string) {
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

// player press up arrow
func ws_up_on_handler(number int, control string) {
	log.Println("ws_up_on_handler. player", number)

	player, ok := game.Players[pn[number]]
	if !ok {
		return
	}

	unpress_other_arrows(player, control)
	player.up_pressed = true

	go ws_up_handler(number, control)
}

func ws_up_handler(number int, control string) {
	log.Println("ws_up_handler. player", number)

	// not beautiful, call it again but player can be disconnected etc, until timer is done
	player, ok := game.Players[pn[number]]
	if !ok {
		return
	}

	//unix time stamp in milliseconds
	unix_ts := time.Now().UnixNano()

	if !player.moving {
		// check if player can move up
		target_position_cell := fmt.Sprintf("%d%d", player.X, player.Y-1)
		if _, ok := game.free_cells[target_position_cell]; !ok {
			return
		}

		seconds := 1
		oneSecond := time.Duration(seconds) * time.Second
		// time in nanoseconds
		one_cell_move_duration := oneSecond.Nanoseconds() / move_cps[player.Turbo]
		player.one_cell_move_duration = one_cell_move_duration
		player.moving_start_time_stamp = unix_ts
		player.Target_y = player.Y - 1
		player.moving = true
		game.Players[pn[number]] = player
		//todo: send to all clients the command to move player. from present position to target position
		time.Sleep(time.Duration(one_cell_move_duration) / 10 * 8 * time.Nanosecond)
		ws_up_handler(number, control)
		return
	} else if player.moving &&
		player.one_cell_move_duration < unix_ts-player.moving_start_time_stamp {
		// check if moving completed
		player.moving = false
		player.Y = player.Target_y
		game.Players[pn[number]] = player
		//todo: send to all clients the command to stand player on target position
		return
	} else if player.moving &&
		player.up_pressed &&
		player.Target_y == player.Y-1 &&
		player.one_cell_move_duration/10*8-100000 < unix_ts-player.moving_start_time_stamp {
		// check if 80% of timer is done, and player still wants to move up
		// if yes, then move/aim player up one cell more if possible
		// 100000 is 100 microseconds gap to wait while after timesleep is over, recursive call is done
		target_position_cell := fmt.Sprintf("%d%d", player.X, player.Y-2)
		if _, ok := game.free_cells[target_position_cell]; !ok {
			time.Sleep(time.Duration(player.one_cell_move_duration) / 100 * time.Nanosecond)
			ws_up_handler(number, control)
			return
		}
		player.Target_y = player.Y - 2
		player.moving_start_time_stamp += player.one_cell_move_duration
		game.Players[pn[number]] = player
		//todo: send to all clients the command to move player. from present position to target position
		time.Sleep(time.Duration(player.one_cell_move_duration) * 120 / 100 * time.Nanosecond)
		ws_up_handler(number, control)
	} else {
		time.Sleep(time.Duration(player.one_cell_move_duration) / 100 * time.Nanosecond)
		ws_up_handler(number, control)
		return
	}

}
