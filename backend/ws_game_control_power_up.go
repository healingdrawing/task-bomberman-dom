package main

import (
	"fmt"
	"log"
)

// "bombs_max++", "explosion_range++", "turbo boolean/move faster" is "1" "2" "3" respectively
func check_power_up(player *PLAYER) {
	// check if player is on power up cell
	// if yes, then apply power up, and hide power up
	// if no, then do nothing

	power_up_to_delete := "" // xy cell indices
	effect := ""             // "1" "2" "3" respectively , to remove css class animation

	game.Power_ups.Range(func(key, value interface{}) bool {
		xy := key.(string)
		power_up := value.(POWER_UP)
		player_xy := fmt.Sprintf("%d%d", player.X, player.Y)
		target_xy := fmt.Sprintf("%d%d", player.Target_x, player.Target_y)
		if power_up.Show && player_xy == xy && target_xy == xy {
			switch power_up.Effect {
			case "1":
				player.bombs_max++
			case "2":
				player.explosion_range++
			case "3":
				player.Turbo = !player.Turbo
			}
			// power_up.Show = false
			if power_up_to_delete == "" {
				power_up_to_delete = xy
				effect = power_up.Effect
			}
			return false
		}
		return true
	})

	if power_up_to_delete != "" {
		game.Power_ups.Delete(power_up_to_delete)
		ws_send_hide_power_up_command(power_up_to_delete, effect)
	}
}

func ws_send_hide_power_up_command(xy string, effect string) {
	log.Println("ws_send_hide_power_up_command")
	message := WS_HIDE_POWER_UP_DTO{
		Cell_xy: xy,
		Effect:  effect,
	}
	uuids := get_all_clients_uuids(clients)
	wsSend(WS_HIDE_POWER_UP, message, uuids)
}

type WS_HIDE_POWER_UP_DTO struct {
	Cell_xy string `json:"cell_xy"`
	Effect  string `json:"effect"`
}
