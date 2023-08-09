package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
)

/*web socket message types. The same as on client(frontend) side*/
type WSMT string

const (
	WS_ERROR_RESPONSE WSMT = "error_response"
	WS_CHAT_MESSAGE   WSMT = "chat_message"

	WS_CLIENT_CONNECTED_TO_SERVER WSMT = "client_connected_to_server"
	WS_CONNECT_TO_SERVER          WSMT = "connect_to_server"
	WS_KEEP_CONNECTION            WSMT = "keep_connection"
	WS_KILL_CONNECTION            WSMT = "kill_connection"
	WS_STILL_CONNECTED            WSMT = "still_connected"
	WS_BROADCAST_MESSAGE          WSMT = "broadcast_message"
	WS_START_GAME                 WSMT = "start_game"
	WS_END_GAME                   WSMT = "end_game"
	WS_PLAYER_GAME_OVER           WSMT = "player_game_over"

	WS_UP            WSMT = "up" // for send to control items on screen
	WS_DOWN          WSMT = "down"
	WS_LEFT          WSMT = "left"
	WS_RIGHT         WSMT = "right"
	WS_STAND         WSMT = "stand"
	WS_BOMB          WSMT = "bomb"
	WS_EXPLODE       WSMT = "explode"
	WS_HIDE_POWER_UP WSMT = "hide_power_up"

	WS_UP_ON     WSMT = "up_on" // for receive from client key press/release events
	WS_UP_OFF    WSMT = "up_off"
	WS_DOWN_ON   WSMT = "down_on"
	WS_DOWN_OFF  WSMT = "down_off"
	WS_LEFT_ON   WSMT = "left_on"
	WS_LEFT_OFF  WSMT = "left_off"
	WS_RIGHT_ON  WSMT = "right_on"
	WS_RIGHT_OFF WSMT = "right_off"
	WS_BOMB_ON   WSMT = "bomb_on"
	WS_BOMB_OFF  WSMT = "bomb_off"

	WS_PLAYER_LIFES      WSMT = "player_lifes"
	WS_CONNECTED_PLAYERS WSMT = "connected_players"
)

type WS_CLIENT_CONNECTED_TO_SERVER_DTO struct {
	Client_number int    `json:"client_number"`
	Client_uuid   string `json:"client_uuid"`
}

// send to client his number and uuid when he connected to server
func ws_client_connected_to_server_handler(client *Client) {
	message := WS_CLIENT_CONNECTED_TO_SERVER_DTO{
		Client_number: client.NUMBER,
		Client_uuid:   client.UUID,
	}
	wsSend(WS_CLIENT_CONNECTED_TO_SERVER, message, []string{client.UUID})
}

type WS_BROADCAST_MESSAGE_DTO struct {
	Content       string `json:"content"`
	Client_number int    `json:"client_number"` // 1..4 or 0 if server or no need coloring
}

func ws_connect_to_server_handler(client *Client, messageData map[string]interface{}) {
	log.Println("=== ws_connect_to_server_handler ===")
	defer wsRecover(messageData)

	uuid, ok := messageData["client_uuid"].(string)
	if !ok {
		log.Println("failed to get client_uuid from messageData")
		return
	}

	nickname, ok := messageData["nickname"].(string)
	if !ok {
		log.Println("failed to get nickname from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{"failed to get nickname from messageData"}, []string{uuid})
		return
	}

	client.NICKNAME = nickname
	html_content := fmt.Sprintf(`
		<div class="color0">
			new star <div class="color%d">[ %s ]</div> raised
		</div>
	`, client.NUMBER, nickname)
	message := WS_BROADCAST_MESSAGE_DTO{
		Content:       html_content,
		Client_number: client.NUMBER,
	}

	uuids := get_all_clients_uuids(clients)

	ws_send_connected_players_number(len(uuids), uuids) // to update players counter on client side

	wsSend(WS_BROADCAST_MESSAGE, message, uuids)

}

func ws_leave_server_handler(client *Client, err error) {
	log.Println("=== ws_leave_server_handler ===")
	defer wsRecover(nil)

	html_content := fmt.Sprintf(`
		<div class="color0">
			old star <div class="color%d">[ %s ]</div> faded
		</div>
	`, client.NUMBER, client.NICKNAME)
	message := WS_BROADCAST_MESSAGE_DTO{
		Content:       html_content,
		Client_number: client.NUMBER,
	}

	uuids := get_all_clients_uuids(clients)

	ws_send_connected_players_number(len(uuids)-1, uuids) // to update players counter on client side
	//minus 1 because rejection will happen automatically, but message sent before that. when this handler is called, error already happened in read incoming messages loop

	wsSend(WS_BROADCAST_MESSAGE, message, uuids)
}

func ws_server_broadcast_handler(text string) {
	log.Println("=== ws_server_broadcast_handler ===")

	html_content := fmt.Sprintf(`<div class="color0"> %s </div>`, text)
	message := WS_BROADCAST_MESSAGE_DTO{
		Content:       html_content,
		Client_number: 0,
	}

	uuids := get_all_clients_uuids(clients)

	wsSend(WS_BROADCAST_MESSAGE, message, uuids)
}

type WS_CONNECTED_PLAYERS_DTO struct {
	Connected_players string `json:"connected_players"`
}

// send to all clients the number of connected clients/players
func ws_send_connected_players_number(n int, uuids []string) {
	log.Println("=== ws_send_connected_players_number ===")

	message := WS_CONNECTED_PLAYERS_DTO{
		Connected_players: fmt.Sprint(n),
	}

	wsSend(WS_CONNECTED_PLAYERS, message, uuids)
}

// todo: add send player lifes later

/*
messageType must be from "ws_utils.go" constants of WSMT type. But go doesn't support enum.
*/
func wsCreateResponseMessage(messageType WSMT, data interface{}) ([]byte, error) {
	response := WS_RESPONSE_MESSAGE_DTO{
		Type: messageType,
		Data: data,
	}

	log.Println("= wsCreateResponseMessage messageType: ", messageType)
	log.Println("= wsCreateResponseMessage data: ", data)
	log.Println("= wsCreateResponseMessage response: ", response)

	jsonData, err := json.Marshal(response)
	if err != nil {
		response.Type = WS_ERROR_RESPONSE
		response.Data = "Error while marshaling response message"
		stableJsonErrorData, _ := json.Marshal(response)
		return stableJsonErrorData, err
	}

	// todo: debug giant print in time of picture sending, so commented
	log.Println("CREATED ================ \nwsCreateResponseMessage: ", string(jsonData))

	return jsonData, nil
}

// wsRecover recover from panic and send a json err response over websocket
func wsRecover(messageData map[string]interface{}) {

	if r := recover(); r != nil {
		fmt.Println("=====================================")
		stackTrace := debug.Stack()
		lines := strings.Split(string(stackTrace), "\n")
		relevantPanicLines := []string{}
		for _, line := range lines {
			if strings.Contains(line, "backend/") {
				relevantPanicLines = append(relevantPanicLines, line)
			}
		}
		if len(relevantPanicLines) > 1 {
			for i, line := range relevantPanicLines {
				if strings.Contains(line, "utils.go") {
					relevantPanicLines = append(relevantPanicLines[:i], relevantPanicLines[i+1:]...)
				}
			}
		}
		relevantPanicLine := strings.Join(relevantPanicLines, "\n")
		log.Println(relevantPanicLines)

		fmt.Println("=====================================")
		// to print the full stack trace
		log.Println(string(stackTrace))

		if messageData == nil {
			log.Println("=== wsRecover:\n=== messageData is nil")
			log.Println("=== emergency close connection")
			return
		}

		uuid, ok := messageData["client_uuid"].(string)
		if !ok {
			log.Println("=== wsRecover: === \n=== failed to get client_uuid from message data")
			return
		}

		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{Content: relevantPanicLine}, []string{uuid})

	}
}
