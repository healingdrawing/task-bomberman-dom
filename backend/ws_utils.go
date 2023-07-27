package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
)

/*websocket message type*/
type WSMT string

const (
	WS_ERROR_RESPONSE     WSMT = "error_response"
	WS_GROUP_CHAT_MESSAGE WSMT = "group_chat_message"

	WS_CONNECT_TO_SERVER          WSMT = "connect_to_server"
	WS_KEEP_CONNECTION            WSMT = "keep_connection"
	WS_KILL_CONNECTION            WSMT = "kill_connection"
	WS_CHAT_MESSAGE               WSMT = "chat_message"
	WS_CLIENT_CONNECTED_TO_SERVER WSMT = "client_connected_to_server"
	WS_STILL_CONNECTED            WSMT = "still_connected"
	WS_BROADCAST_MESSAGE          WSMT = "broadcast_message"
	WS_START_GAME                 WSMT = "start_game"
	WS_END_GAME                   WSMT = "end_game"
	WS_PLAYER_GAME_OVER           WSMT = "player_game_over"
	WS_UP                         WSMT = "up"
	WS_DOWN                       WSMT = "down"
	WS_LEFT                       WSMT = "left"
	WS_RIGHT                      WSMT = "right"
	WS_STAND                      WSMT = "stand"
	WS_BOMB                       WSMT = "bomb"
	WS_EXPLODE                    WSMT = "explode"
	WS_HIDE_POWER_UP              WSMT = "hide_power_up"
)

/*
messageType must be from "ws_utils.go" constants of WSMT type. But go doesn't support enum.
*/
func wsCreateResponseMessage(messageType WSMT, data interface{}) ([]byte, error) {
	response := WS_RESPONSE_MESSAGE_DTO{
		Type: messageType,
		Data: data,
	}

	log.Println("= wsCreateResponseMessage: ", messageType)

	jsonData, err := json.Marshal(response)
	if err != nil {
		response.Type = WS_ERROR_RESPONSE
		response.Data = "Error while marshaling response message"
		stableJsonErrorData, _ := json.Marshal(response)
		return stableJsonErrorData, err
	}

	// todo: debug giant print in time of picture sending, so commented
	// log.Println("CREATED ================ \nwsCreateResponseMessage: ", string(jsonData))

	return jsonData, nil
}

// wsRecover recover from panic and send a json err response over websocket
func wsRecover(messageData map[string]interface{}) {

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("=== wsRecover: === \n=== failed to get user_uuid from message data")
		return
	}

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

		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{Content: relevantPanicLine}, []string{uuid})
		fmt.Println("=====================================")
		// to print the full stack trace
		log.Println(string(stackTrace))
	}
}
