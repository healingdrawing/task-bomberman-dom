package main

import (
	"log"
	"strconv"
)

// fill initial position and dead property for each player
func prepare_players() {
	game.Players = make(map[string]PLAYER)

	// fill x4 players, because it is the same data for all players, later the dead property will be changed, according to Client.NUMBER
	for i := 1; i < 5; i++ {
		game.Players[strconv.Itoa(i)] = PLAYER{
			Number:          i,
			X:               6,
			Y:               0,
			Target_x:        6,
			Target_y:        0,
			bombs_max:       1,
			bombs_left:      1,
			explosion_range: 1,
			Turbo:           false,
			Dead:            true,
		}
	}

	// set the starting positions for each player
	px := []int{0, 6, 0, 6}
	py := []int{0, 6, 6, 0}

	pn := map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}

	// iterate over all clients and switch their dead property to false, also set their starting position as x and y, and target_x and target_y
	clients.Range(func(key, value interface{}) bool {
		number := value.(*Client).NUMBER
		log.Println("number", number)
		player, ok := game.Players[pn[number]]
		if !ok {
			log.Fatalln("cant find player by number", number, "in game.Players", game.Players)
		}
		player.Dead = false
		player.X = px[number-1]
		player.Y = py[number-1]
		player.Target_x = px[number-1]
		player.Target_y = py[number-1]
		game.Players[pn[number]] = player //todo: without this line, the player is not updated in game.Players
		return true
	})
}

// fill initial position for each weak obstacle
func prepare_weak_obstacles_and_power_ups() {

	game.Weak_obstacles = make(map[string]bool)

	// yes it is only x4 of them(not a best time for experimenting)
	var x, y []int
	external := randomNum(0, 1) > 0 // external x4 or internal x4 weak obstacles
	if external {
		x = []int{6, 3, 0, 3}
		y = []int{3, 0, 3, 6}
	} else {
		x = []int{4, 3, 2, 3}
		y = []int{3, 2, 3, 4}
	}

	// the key is the position in same time
	for i := 0; i < 4; i++ {
		game.Weak_obstacles[strconv.Itoa(x[i])+strconv.Itoa(y[i])] = true
	}

	// fill the powerups
	// shuffle the powerups
	pups := []string{"bombs_max", "explosion_range", "turbo", "turbo"}
	for i := len(pups) - 1; i > 0; i-- {
		j := randomNum(0, i)
		pups[i], pups[j] = pups[j], pups[i]
	}

	game.Power_ups = make(map[string]POWER_UP)
	for i := 0; i < 4; i++ {
		game.Power_ups[strconv.Itoa(x[i])+strconv.Itoa(y[i])] = POWER_UP{
			Effect: pups[i],
			Show:   false,
		}
	}
}
