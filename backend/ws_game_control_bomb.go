package main

import (
	"fmt"
	"log"
	"time"
)

func ws_bomb_off_handler(number string) {
	dprint("ws_bomb_off_handler. player", number)

	player_value, ok := game.Players.Load(number)
	if !ok {
		return
	}
	player := player_value.(PLAYER)
	player.bomb_pressed = false
	game.Players.Store(number, player)
}

// player press Enter, bomb will be placed
func ws_bomb_handler(number string) {
	dprint("ws_bomb_handler. player", number)
	player_value, ok := game.Players.Load(number)
	if !ok {
		return
	}

	player := player_value.(PLAYER)

	// Unix timestamp in nanoseconds
	unix_ts := time.Now().UnixNano()

	if player.bombs_used >= player.bombs_max {
		return
	}

	// check if bomb is alredy placed on player cell, then do not place another bomb
	bomb_xy := fmt.Sprintf("%d%d", player.X, player.Y)
	if _, ok := game.bomb_cells.Load(bomb_xy); ok {
		return
	}

	// create the bomb
	player.bombs_used++
	game.Players.Store(number, player)
	game.bomb_cells.Store(bomb_xy, unix_ts)

	// send bomb to all players
	ws_send_bomb_command(player.Number, bomb_xy)

	// bomb will explode after (player.explosion_range +2) seconds, to escape
	go func() {
		explosion_range := player.explosion_range
		time.Sleep(time.Duration(explosion_range+2) * time.Second)
		ws_explosion_handler(player.Number, bomb_xy, explosion_range)
	}()

}

func ws_send_bomb_command(player_number int, cell_xy string) {
	log.Println("ws_send_bomb_command. player", player_number)
	message := WS_BOMB_DTO{
		Number:    player_number,
		Target_xy: cell_xy,
	}
	uuids := get_all_clients_uuids(clients)
	wsSend(WS_BOMB, message, uuids)
}

type WS_BOMB_DTO struct {
	Number    int    `json:"number"`
	Target_xy string `json:"target_xy"`
}

// bomb handler call this as goroutine
func ws_explosion_handler(player_number int, bomb_xy string, explosion_range int) {
	dprint("ws_explosion_handler. player", player_number, "bomb_xy", bomb_xy, "explosion_range", explosion_range)
	// remove bomb from bomb_cells
	game.bomb_cells.Delete(bomb_xy)

	// decrease player.bombs_used, if player is still alive/connected. if something wrong, just ignore it
	player_value, ok := game.Players.Load(string_number[player_number])
	if ok {
		player := player_value.(PLAYER)
		player.bombs_used--
		game.Players.Store(string_number[player_number], player)
	} else {
		dprint("ws_explosion_handler. player", player_number, "not found")
	}

	rs := []rune(bomb_xy)

	if len(rs) != 2 {
		return
	}

	bomb_x, ok := cell_number_from_string[string(rs[0])]
	if !ok {
		return
	}
	bomb_y, ok := cell_number_from_string[string(rs[1])]
	if !ok {
		return
	}

	explosion_cells_xy := []string{bomb_xy}
	// calculate explosion cells in all directions, based on explosion_range, from bomb_xy

	for i := 1; i <= explosion_range; i++ {
		// up
		if bomb_y-i > -1 { // 0 is first cell
			explosion_cells_xy = append(explosion_cells_xy, fmt.Sprintf("%d%d", bomb_x, bomb_y-i))
		}
		// down
		if bomb_y+i < 7 { // 6 is last cell
			explosion_cells_xy = append(explosion_cells_xy, fmt.Sprintf("%d%d", bomb_x, bomb_y+i))
		}
		// left
		if bomb_x-i > -1 { // 0 is first cell
			explosion_cells_xy = append(explosion_cells_xy, fmt.Sprintf("%d%d", bomb_x-i, bomb_y))
		}
		// right
		if bomb_x+i < 7 { // 6 is last cell
			explosion_cells_xy = append(explosion_cells_xy, fmt.Sprintf("%d%d", bomb_x+i, bomb_y))
		}

	}

	//todo: implement affecting players and weak obstacles
	// 1- check the players cells, if player is on explosion or aimed to explosion, then if player has lifes, decrease the lifes of the player else kill the player

	game.Players.Range(func(key, value interface{}) bool {
		player := value.(PLAYER)
		player_xy := fmt.Sprintf("%d%d", player.X, player.Y)
		target_xy := fmt.Sprintf("%d%d", player.Target_x, player.Target_y)
		for _, explosion_cell_xy := range explosion_cells_xy {
			if player_xy == explosion_cell_xy || target_xy == explosion_cell_xy {
				player.Lifes--
				if player.Lifes > 0 {
					game.Players.Store(key, player)
				} else {
					game.Players.Delete(key) // remove from range loop etc, let is say, player is dead
				}
				ws_send_player_lifes(player.Number, player.Lifes, player.uuid)

			}
		}
		return true
	})

	// 2- check the weak obstacles cells, if weak obstacle is on explosion, then remove the weak obstacle, add free cell, and show the powerup on the cell(because each of the four obstacles have a powerup)

	destroyed_weak_obstacles := []string{}
	appeared_power_up_effect := []string{}

	game.Weak_obstacles.Range(func(key, value interface{}) bool {
		weak_obstacle_xy := key.(string)
		for _, explosion_cell_xy := range explosion_cells_xy {
			if weak_obstacle_xy == explosion_cell_xy {
				destroyed_weak_obstacles = append(destroyed_weak_obstacles, weak_obstacle_xy)
				game.Weak_obstacles.Delete(key)
				// add free cell
				game.free_cells.Store(weak_obstacle_xy, true)
				// show the powerup on the cell
				power_up_data, ok := game.Power_ups.Load(weak_obstacle_xy)
				if ok {
					power_up := power_up_data.(POWER_UP)
					power_up.Show = true
					appeared_power_up_effect = append(appeared_power_up_effect, power_up.Effect)
					game.Power_ups.Store(weak_obstacle_xy, power_up)
				} else {
					dprint("ws_explosion_handler. power_up not found for weak_obstacle_xy", weak_obstacle_xy)
				}
				break
			}
		}
		return true
	})

	// 3- extend response object for explosion, and send to all clients

	// send explosion to all players, must manage also all affected items: players, weak obstacles. First remove all affected items, then execute explosion on client side
	ws_send_explosion_command(player_number, bomb_xy, explosion_cells_xy, destroyed_weak_obstacles, appeared_power_up_effect)

}

func ws_send_player_lifes(player_number int, player_lifes int, player_uuid string) {
	log.Println("ws_send_player_lifes. player", player_number, "lifes", player_lifes)
	message := WS_PLAYER_LIFES_DTO{
		Number: player_number,
		Lifes:  player_lifes,
	}

	wsSend(WS_PLAYER_LIFES, message, []string{player_uuid})
}

type WS_PLAYER_LIFES_DTO struct {
	Number int `json:"number"` // to remove bomb animation for player who placed the bomb
	Lifes  int `json:"lifes"`  // the first one is bomb_xy, to remove bomb
}

// todo: extend this and function above and send also affected items commands
func ws_send_explosion_command(player_number int, bomb_xy string, explosion_cells_xy []string, destroyed_weak_obstacles []string, appeared_power_up_effect []string) {
	log.Println("ws_send_explosion_command. bomb_xy", bomb_xy, "explosion_cells_xy", explosion_cells_xy)
	message := WS_EXPLODE_DTO{
		Number:          player_number,
		Cells_xy:        explosion_cells_xy,
		Destroy_xy:      destroyed_weak_obstacles,
		Power_up_effect: appeared_power_up_effect,
	}
	uuids := get_all_clients_uuids(clients)
	wsSend(WS_EXPLODE, message, uuids)
}

type WS_EXPLODE_DTO struct {
	Number          int      `json:"number"`          // to remove bomb animation for player who placed the bomb
	Cells_xy        []string `json:"cells_xy"`        // the first one is bomb_xy, to remove bomb
	Destroy_xy      []string `json:"destroy_xy"`      // xy to destroy weak obstacles
	Power_up_effect []string `json:"power_up_effect"` // power up effect to replace weak obstacle
}
