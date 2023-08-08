package main

import (
	"fmt"
	"log"
	"time"
)

func ws_right_off_handler(number string, control string) {
	dprint("ws_right_off_handler. player", number)

	player_value, ok := game.Players.Load(number)
	if !ok {
		return
	}
	player := player_value.(PLAYER)
	unpress_arrow(&player, control)
	game.Players.Store(number, player)
}

// player press right arrow
func ws_right_handler(number string, control string, press bool) {
	player_value, ok := game.Players.Load(number)
	if !ok {
		return
	}

	player := player_value.(PLAYER)

	if player.direction != RIGHT && player.direction != STAND {
		return
	}

	if press {
		press_arrow_unpress_other_arrows(&player, control)
		game.Players.Store(number, player)
	}

	// Unix timestamp in nanoseconds
	unix_ts := time.Now().UnixNano()

	if !player.moving && press {
		dprint("inside if !player.moving && press")
		// Check if player can move right
		target_position_cell := fmt.Sprintf("%d%d", player.X+1, player.Y)
		if _, ok := game.free_cells.Load(target_position_cell); !ok {
			dprint("target cell is not free", target_position_cell)
			player.right_pressed = false
			return
		}
		seconds := 1
		one_second := time.Duration(seconds) * time.Second
		// Time in nanoseconds
		one_cell_move_duration := one_second.Nanoseconds() / move_cps[player.Turbo]
		dprint("one_cell_move_duration", one_cell_move_duration)
		player.one_cell_move_duration = one_cell_move_duration
		player.moving_start_time_stamp = unix_ts
		player.Target_x = player.X + 1
		player.can_change_present_cell = true
		player.can_change_target_cell = false
		player.moving = true
		player.direction = RIGHT
		game.Players.Store(number, player)
		ws_send_move_right_command(&player)
	} else if player.moving && player.can_change_target_cell &&
		player.right_pressed &&
		player.one_cell_move_duration/10*8 < unix_ts-player.moving_start_time_stamp {
		// Check if 80% of the timer is done, and the player still wants to move right
		// If yes, then move/aim the player right one cell more if possible
		dprint("inside if time 80")
		target_position_cell := fmt.Sprintf("%d%d", player.X+1, player.Y)
		fc, ok := game.free_cells.Load(target_position_cell)
		if !ok {
			dprint("target cell is not free", target_position_cell)
			player.right_pressed = false
			return
		}
		dprint("free cell", fc)

		player.Target_x = player.X + 1
		player.can_change_present_cell = true
		player.can_change_target_cell = false
		player.moving_start_time_stamp += player.one_cell_move_duration
		game.Players.Store(number, player)
		dprint("inside if time 80. before ws_send_move_right_command")
		ws_send_move_right_command(&player)
	} else if player.moving &&
		player.can_change_present_cell &&
		player.one_cell_move_duration/10*5 < unix_ts-player.moving_start_time_stamp {
		// Check if 50% of the timer is done, switch to the next cell
		dprint("inside if time 50")
		player.X = player.Target_x
		player.can_change_present_cell = false
		player.can_change_target_cell = true
		game.Players.Store(number, player)
	} else if player.moving &&
		player.one_cell_move_duration < unix_ts-player.moving_start_time_stamp {
		dprint("inside if time 100")
		player.can_change_present_cell = false
		player.can_change_target_cell = false
		player.moving = false
		player.direction = STAND
		game.Players.Store(number, player)
	}
}

func ws_send_move_right_command(player *PLAYER) {
	log.Println("ws_send_move_right_command. player", player.Number)
	message := WS_MOVE_RIGHT_DTO{
		Number:   player.Number,
		Target_x: player.Target_x,
		Turbo:    player.Turbo,
	}
	dprint("ws_send_move_right_command. message", message)
	uuids := get_all_clients_uuids(clients)
	dprint("ws_send_move_right_command. uuids", uuids)
	wsSend(WS_RIGHT, message, uuids)
	dprint("ws_send_move_right_command. sent")

}

type WS_MOVE_RIGHT_DTO struct {
	Number   int  `json:"number"`
	Target_x int  `json:"target_x"`
	Turbo    bool `json:"turbo"`
}
