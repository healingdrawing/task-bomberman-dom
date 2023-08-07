package main

// player press up arrow
// func ws_up_on_handler(number int, control string) {
// 	log.Println("ws_up_on_handler. player", number)

// 	player, ok := game.Players[pn[number]]
// 	if !ok {
// 		return
// 	}

// 	unpress_other_arrows(&player, control)
// 	player.up_pressed = true
// 	game.Players[pn[number]] = player

// 	go ws_up_handler(number, control)
// }

// func ws_up_handler(number int, control string) {
// 	log.Println("ws_up_handler. player", number)

// 	// not beautiful, call it again but player can be disconnected etc, until timer is done
// 	player, ok := game.Players[pn[number]]
// 	if !ok {
// 		return
// 	}

// 	//unix time stamp in milliseconds
// 	unix_ts := time.Now().UnixNano()

// 	if !player.moving {
// 		// check if player can move up
// 		target_position_cell := fmt.Sprintf("%d%d", player.X, player.Y-1)
// 		if _, ok := game.free_cells[target_position_cell]; !ok {
// 			return
// 		}

// 		seconds := 1
// 		oneSecond := time.Duration(seconds) * time.Second
// 		// time in nanoseconds
// 		one_cell_move_duration := oneSecond.Nanoseconds() / move_cps[player.Turbo]
// 		player.one_cell_move_duration = one_cell_move_duration
// 		player.moving_start_time_stamp = unix_ts
// 		player.Target_y = player.Y - 1
// 		player.moving = true
// 		game.Players[pn[number]] = player
// 		//todo: send to all clients the command to move player. from present position to target position
// 		ws_send_move_up_command(&player)
// 		time.Sleep(time.Duration(one_cell_move_duration) / 10 * 8 * time.Nanosecond)
// 		ws_up_handler(number, control)
// 		return
// 	} else if player.moving &&
// 		player.one_cell_move_duration < unix_ts-player.moving_start_time_stamp {
// 		// check if moving completed
// 		player.moving = false
// 		player.Y = player.Target_y
// 		unpress_all_arrows(&player)
// 		game.Players[pn[number]] = player
// 		//todo: send to all clients the command to stand player on target position. not sure this needed
// 		return
// 	} else if player.moving &&
// 		player.up_pressed &&
// 		player.Target_y == player.Y-1 &&
// 		player.one_cell_move_duration/10*8-100000 < unix_ts-player.moving_start_time_stamp {
// 		// check if 80% of timer is done, and player still wants to move up
// 		// if yes, then move/aim player up one cell more if possible
// 		// 100000 is 100 microseconds gap to wait while after timesleep is over, recursive call is done
// 		target_position_cell := fmt.Sprintf("%d%d", player.X, player.Y-2)
// 		if _, ok := game.free_cells[target_position_cell]; !ok {
// 			time.Sleep(time.Duration(player.one_cell_move_duration) / 100 * time.Nanosecond)
// 			ws_up_handler(number, control)
// 			return
// 		}
// 		player.Target_y = player.Y - 2
// 		player.moving_start_time_stamp += player.one_cell_move_duration
// 		game.Players[pn[number]] = player
// 		//todo: send to all clients the command to move player. from present position to target position
// 		ws_send_move_up_command(&player)
// 		time.Sleep(time.Duration(player.one_cell_move_duration) * 120 / 100 * time.Nanosecond)
// 		ws_up_handler(number, control)
// 	} else {
// 		time.Sleep(time.Duration(player.one_cell_move_duration) / 100 * time.Nanosecond)
// 		ws_up_handler(number, control)
// 		return
// 	}

// }
