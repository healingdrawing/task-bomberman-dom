package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	min_players    = 2
	max_players    = 4
	waiting_time   = 20 * time.Second
	countdown_time = 10 * time.Second
)

// WaitingState represents the state of the waiting period and countdown.
type WaitingState int

const (
	WAITING_FOR_PLAYERS WaitingState = iota
	WAITING_FOR_COUNTDOWN
	GAME_STARTED
)

var (
	connected_players_number int
	game_waiting_state       WaitingState
	waitingTimerStarted      bool
	waitingTimerMutex        sync.Mutex
)

// game_waiting_state_handle_client_connected is called when a new client is connected.
// It handles the logic for the waiting and countdown.
func game_waiting_state_handle_client_connected() {
	connected_players_number++

	time.Sleep(time.Second) // gap to print 20 seconds letf, after new star raised

	switch game_waiting_state {
	case WAITING_FOR_PLAYERS:
		if connected_players_number >= min_players {
			waitingTimerMutex.Lock()
			if !waitingTimerStarted {
				waitingTimerStarted = true
				waitingTimerMutex.Unlock()

				fmt.Println("Waiting for more players...")
				ws_server_broadcast_handler("Waiting for more players...")
				for players_countdown := waiting_time / time.Second; players_countdown > 0; players_countdown-- {
					if connected_players_number >= max_players {
						break
					}
					if connected_players_number < min_players {
						game_waiting_state = WAITING_FOR_PLAYERS
						fmt.Println("Waiting countdown canceled, waiting for players.")
						ws_server_broadcast_handler("Waiting countdown canceled, waiting for players.")
						break
					}

					fmt.Printf("%d seconds left\n", players_countdown)
					// Here, you can add the logic to notify clients about the remaining seconds.
					ws_server_broadcast_handler(fmt.Sprintf("%d seconds left", players_countdown))
					time.Sleep(time.Second)
				}

				if connected_players_number >= min_players && connected_players_number <= max_players {
					game_waiting_state = WAITING_FOR_COUNTDOWN
					fmt.Println("Countdown started!")
					ws_server_broadcast_handler("Countdown started!")
					for prepare_countdown := countdown_time / time.Second; prepare_countdown > 0; prepare_countdown-- {

						if connected_players_number < min_players {
							game_waiting_state = WAITING_FOR_PLAYERS
							fmt.Println("Prepare countdown canceled, waiting for players.")
							ws_server_broadcast_handler("Prepare countdown canceled, waiting for players.")
							break
						}

						fmt.Printf("%d seconds left\n", prepare_countdown)
						// Here, you can add the logic to notify clients about the remaining seconds.
						ws_server_broadcast_handler(fmt.Sprintf("%d seconds left", prepare_countdown))
						time.Sleep(time.Second)
					}

					if connected_players_number >= min_players && connected_players_number <= max_players {
						game_waiting_state = GAME_STARTED
						fmt.Println("Game started!")
						// Here, you can add the logic to notify clients that the game has started.
						ws_server_broadcast_handler("!!!GO GO GO!!!")
					} else if connected_players_number < min_players {
						// todo: not sure it can fires , after injection/duplication above
						game_waiting_state = WAITING_FOR_PLAYERS
						fmt.Println("POTENTIAL ARTEFACT Waiting countdown canceled, waiting for players.")
						// Here, you can add the logic to notify clients that the countdown has been canceled.
						ws_server_broadcast_handler("POTENTIAL ARTEFACT Waiting countdown canceled, waiting for players.")
					} else {
						game_waiting_state = WAITING_FOR_PLAYERS
						fmt.Println("Waiting countdown canceled, too many players.")
						// Here, you can add the logic to notify clients that the countdown has been canceled.
						ws_server_broadcast_handler("Waiting countdown canceled, too many players. Unexpected condition. Press \"R.I.P\" to reconnect.")
					}
				}
			} else {
				waitingTimerMutex.Unlock()
			}
		}

	}
}

// game_waiting_state_handle_client_disconnected is called when a client disconnects.
// It handles the logic for resetting the waiting and countdown states.
func game_waiting_state_handle_client_disconnected() {
	connected_players_number--

	switch game_waiting_state {
	case WAITING_FOR_PLAYERS:
		if connected_players_number < min_players {
			waitingTimerMutex.Lock()
			waitingTimerStarted = false
			waitingTimerMutex.Unlock()

			fmt.Println("DISCONNECTED Waiting for more players...") // TODO: remove
		}
	case WAITING_FOR_COUNTDOWN:
		if connected_players_number < min_players {
			waitingTimerMutex.Lock()
			waitingTimerStarted = false
			waitingTimerMutex.Unlock()

			game_waiting_state = WAITING_FOR_PLAYERS
			// Here, you can add the logic to notify clients that the countdown has been canceled.
			fmt.Println("DISCONNECTED prepare for game...") // TODO: remove
		}
	case GAME_STARTED:
		log.Println("DOES NOT RESET TO WAITING FOR PLAYERS, BECAUSE GAME IS ALREADY STARTED")
		// Implement logic for handling a client disconnect during the game, if needed.
	}
}
