package main

import "log"

func ws_character_control_handler(client *Client, control string) {

	switch control {
	case string(WS_UP_ON):
		log.Println("WS_UP_ON", client.NUMBER)
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
