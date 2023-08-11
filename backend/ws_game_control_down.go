package main

import (
	"fmt"
	"log"
	"time"
)

func ws_down_off_handler(number string, control string) {
	dprint("ws_down_off_handler. player", number)

	player_value, ok := game.Players.Load(number)
	if !ok {
		return
	}
	player := player_value.(PLAYER)
	unpress_arrow(&player, control)
	game.Players.Store(number, player)
}

// player press down arrow
func ws_down_handler(number string, control string, press bool) {
	player_value, ok := game.Players.Load(number)
	if !ok {
		return
	}

	player := player_value.(PLAYER)

	if player.direction != DOWN && player.direction != STAND {
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
		// Check if player can move down
		target_position_cell := fmt.Sprintf("%d%d", player.X, player.Y+1)
		if _, ok := game.free_cells.Load(target_position_cell); !ok {
			dprint("target cell is not free", target_position_cell)
			player.down_pressed = false
			return
		}
		seconds := 1
		one_second := time.Duration(seconds) * time.Second
		// Time in nanoseconds
		one_cell_move_duration := one_second.Nanoseconds() / move_cps[player.Turbo]
		dprint("one_cell_move_duration", one_cell_move_duration)
		player.one_cell_move_duration = one_cell_move_duration
		player.moving_start_time_stamp = unix_ts
		player.Target_y = player.Y + 1
		player.can_change_present_cell = true
		player.can_change_target_cell = false
		player.moving = true
		player.direction = DOWN
		game.Players.Store(number, player)
		ws_send_move_down_command(&player)
	} else if player.moving && player.can_change_target_cell &&
		player.down_pressed &&
		player.one_cell_move_duration/10*8 < unix_ts-player.moving_start_time_stamp {
		// Check if 80% of the timer is done, and the player still wants to move down
		// If yes, then move/aim the player down one cell more if possible
		dprint("inside if time 80")
		target_position_cell := fmt.Sprintf("%d%d", player.X, player.Y+1)
		fc, ok := game.free_cells.Load(target_position_cell)
		if !ok {
			dprint("target cell is not free", target_position_cell)
			player.down_pressed = false
			return
		}
		dprint("free cell", fc)

		player.Target_y = player.Y + 1
		player.can_change_present_cell = true
		player.can_change_target_cell = false
		player.moving_start_time_stamp += player.one_cell_move_duration
		game.Players.Store(number, player)
		dprint("inside if time 80. before ws_send_move_down_command")
		ws_send_move_down_command(&player)
	} else if player.moving &&
		player.can_change_present_cell &&
		player.one_cell_move_duration/10*5 < unix_ts-player.moving_start_time_stamp {
		// Check if 50% of the timer is done, switch to the next cell
		dprint("inside if time 50")
		player.Y = player.Target_y
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
		check_power_up(&player)
		game.Players.Store(number, player)
	}
}

func ws_send_move_down_command(player *PLAYER) {
	log.Println("ws_send_move_down_command. player", player.Number)
	message := WS_MOVE_DOWN_DTO{
		Number:   player.Number,
		Target_y: player.Target_y,
		Turbo:    player.Turbo,
	}
	uuids := get_all_clients_uuids(clients)
	wsSend(WS_DOWN, message, uuids)
}

type WS_MOVE_DOWN_DTO struct {
	Number   int  `json:"number"`
	Target_y int  `json:"target_y"`
	Turbo    bool `json:"turbo"`
}
