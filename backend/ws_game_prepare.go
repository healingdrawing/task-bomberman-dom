package main

import "strconv"

// fill initial position and dead property for each player
func prepare_players() {
	game.players = make(map[string]PLAYER)

	// fill x4 players, because it is the same data for all players, later the dead property will be changed, according to Client.NUMBER
	for i := 1; i < 5; i++ {
		game.players[strconv.Itoa(i)] = PLAYER{
			number:          i,
			x:               6,
			y:               0,
			target_x:        6,
			target_y:        0,
			bombs_max:       1,
			bombs_left:      1,
			explosion_range: 1,
			turbo:           false,
			dead:            true,
		}
	}

	// set the starting positions for each player
	px := []int{0, 6, 0, 6}
	py := []int{0, 6, 6, 0}

	// iterate over all clients and switch their dead property to false, also set their starting position as x and y, and target_x and target_y
	clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		player := game.players[strconv.Itoa(client.NUMBER)]
		player.dead = false
		player.x = px[client.NUMBER-1]
		player.y = py[client.NUMBER-1]
		return true
	})
}

// fill initial position for each weak obstacle
func prepare_weak_obstacles() {

	game.weak_obstacles = make(map[string]WEAK_OBSTACLE)

	// yes it is only x4 of them(not a best time for experimenting)
	// internal := randomNum(0, 1) > 0 // external x4 or internal x4 weak obstacles
	// todo: continue, time to switch machine

}
