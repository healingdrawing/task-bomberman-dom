package main

import (
	"log"
	"strconv"
	"sync"
)

// fill initial position and dead property for each player
func prepare_players() {
	game.Players = sync.Map{} // Initialize Players as a sync.Map

	// Fill x4 players, because it is the same data for all players
	for i := 1; i < 5; i++ {
		number := string_number[i]
		player := PLAYER{
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
			up_pressed:      false,
			down_pressed:    false,
			left_pressed:    false,
			right_pressed:   false,
			bomb_pressed:    false,
		}
		game.Players.Store(number, player)
	}

	// Set the starting positions for each player
	px := []int{0, 6, 0, 6}
	py := []int{0, 6, 6, 0}

	// Iterate over all clients and update player properties in the sync.Map
	clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		number := client.NUMBER
		log.Println("number", number)
		number_str := string_number[number]
		if player, ok := game.Players.Load(number_str); ok {
			player_data := player.(PLAYER)
			player_data.Dead = false
			player_data.X = px[number-1]
			player_data.Y = py[number-1]
			player_data.Target_x = px[number-1]
			player_data.Target_y = py[number-1]
			game.Players.Store(number_str, player_data)
		} else {
			log.Fatalln("can't find player by number", number, "in game.Players")
		}
		return true
	})
}

// fill initial position for each weak obstacle
func prepare_weak_obstacles_and_power_ups() {
	game.Weak_obstacles = sync.Map{} // Initialize Weak_obstacles as a sync.Map

	var x, y []int
	external := randomNum(0, 1) > 0 // external x4 or internal x4 weak obstacles
	if external {
		x = []int{6, 3, 0, 3}
		y = []int{3, 0, 3, 6}
	} else {
		x = []int{4, 3, 2, 3}
		y = []int{3, 2, 3, 4}
	}

	// The key is the position in the form "xy"
	for i := 0; i < 4; i++ {
		game.Weak_obstacles.Store(strconv.Itoa(x[i])+strconv.Itoa(y[i]), true)
	}

	// Fill the power-ups
	// Shuffle the power-ups
	pups := []string{"bombs_max", "explosion_range", "turbo", "turbo"}
	for i := len(pups) - 1; i > 0; i-- {
		j := randomNum(0, i)
		pups[i], pups[j] = pups[j], pups[i]
	}

	game.Power_ups = sync.Map{} // Initialize Power_ups as a sync.Map
	for i := 0; i < 4; i++ {
		game.Power_ups.Store(strconv.Itoa(x[i])+strconv.Itoa(y[i]), POWER_UP{
			Effect: pups[i],
			Show:   false,
		})
	}
}

// fill the locked for moving cells xy, it is strong obstacles and weak obstacles(not destroyed yet)
func prepare_free_cells() {
	game.free_cells = sync.Map{} // Initialize free_cells as a sync.Map
	locked_cells := sync.Map{}   // Initialize locked_cells as a sync.Map

	// Fill the strong obstacles x9, every second from the left top corner
	strong_x := []int{1, 3, 5, 1, 3, 5, 1, 3, 5}
	strong_y := []int{1, 1, 1, 3, 3, 3, 5, 5, 5}
	for i := 0; i < 9; i++ {
		locked_cells.Store(strconv.Itoa(strong_x[i])+strconv.Itoa(strong_y[i]), true)
	}

	// Add the weak obstacles to the locked_cells
	game.Weak_obstacles.Range(func(key, value interface{}) bool {
		k := key.(string)
		locked_cells.Store(k, true)
		return true
	})

	// Fill the free cells
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			xy := strconv.Itoa(i) + strconv.Itoa(j)
			if _, ok := locked_cells.Load(xy); !ok {
				game.free_cells.Store(xy, true)
			}
		}
	}
}
