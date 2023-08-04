package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	clients  = &sync.Map{}
)

type wsInput struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func wsConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("=== inside wsConnection ===")

	if connected_clients_number(clients) > 3 {
		log.Println("=== too many clients ===")
		jsonResponse(w, http.StatusOK, "Too many clients")
		return
	}

	number := choose_first_free_number(clients)
	if number == 0 {
		log.Println("=== no free numbers ===")
		jsonResponse(w, http.StatusOK, "No free slots")
		return
	}

	uuid, err := generate_UUID()
	log.Println("wsConnection uuid: ", uuid) //todo: delete debug
	if err != nil {
		log.Println("====================================")
		log.Println("uuid creation error: ", err)
		log.Println("====================================")
		jsonResponse(w, http.StatusOK, "uuid creation error. status 200 , because otherwise no message in browser console, facepalm")
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	reader(ws, number, uuid)
}

type Client struct {
	CONN     *websocket.Conn
	NUMBER   int
	UUID     string
	NICKNAME string
}

func reader(conn *websocket.Conn, client_number int, uuid string) {
	log.Println("=== inside reader ===")

	client := &Client{CONN: conn, NUMBER: client_number, UUID: uuid}
	clients.Store(uuid, client)
	defer clients.Delete(uuid) //todo: looks like this needed, forgot when i commented this and why
	defer conn.Close()

	ws_client_connected_to_server_handler(client)
	go game_waiting_state_handle_client_connected()

	for {
		messageType, incoming, err := conn.ReadMessage()
		log.Println("=== inside reader ===")
		if err != nil {
			log.Println("=== error in reader ===")
			log.Println("messageType, incoming, err := conn.ReadMessage()")
			log.Println("messageType", messageType)
			log.Println("incoming", incoming)
			log.Println("err", err)
			log.Println("=== error in reader , before delete and close ws ===")
			ws_leave_server_handler(client, err)
			game_waiting_state_handle_client_disconnected()
			return
		}

		// todo: debug print
		log.Println("=================\nread message:", "\nincoming as string:", string(incoming), "\nmessageType: ", messageType) //todo: delete debug

		if messageType == websocket.TextMessage {
			log.Println("Text message received")
			var data wsInput
			if err := json.Unmarshal(incoming, &data); err != nil {
				log.Println(err)
				return
			}

			// todo: debug print
			log.Println("data after unmarshalling: ", data) //todo: delete debug

			switch data.Type {
			case string(WS_CONNECT_TO_SERVER):
				log.Println("==================WS_CONNECT_TO_SERVER FIRED==================")
				log.Println("data.Data: ", data.Data)
				ws_connect_to_server_handler(client, data.Data)

			case string(WS_CHAT_MESSAGE):
				wsChatMessageHandler(conn, data.Data)

			default:
				log.Println("Unknown type: ", data.Type)
			}
		}
		if messageType == websocket.BinaryMessage {
			log.Println("Binary message received")
		}
	}
}

// send message to connections by uuids provided
func wsSend(message_type WSMT, message interface{}, uuids []string) {
	outputMessage, err := wsCreateResponseMessage(message_type, message)

	if err != nil {
		log.Println(err)
	}

	for _, uuid := range uuids {
		if raw_client, ok := clients.Load(uuid); ok {
			if client, ok := raw_client.(*Client); ok {
				err = client.CONN.WriteMessage(websocket.TextMessage, outputMessage)
				log.Println("wsSend: message sent to client: ", uuid)
				if err != nil {
					log.Println(err)
				}
			} else {
				log.Println("wsSend: clients.Load(uuid) is not a *Client")
			}
		} else {
			log.Println("wsSend: client not found . clients.Load(uuid) failed")
		}
	}
}
