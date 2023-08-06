package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WS_CHAT_MESSAGE_DTO struct {
	Content       string `json:"content"`
	Nickname      string `json:"nickname"`
	Client_number int    `json:"client_number"`
	Created_at    string `json:"created_at"`
}

// todo: can be refactored to use client instead of conn. with comparing of uuids to verify that it is the same client
func wsChatMessageHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["client_uuid"].(string)
	if !ok {
		log.Println("failed to get client_uuid from messageData")
		return
	}

	content, ok := messageData["content"].(string)
	if !ok {
		log.Println("failed to get content from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " failed to get content from messageData")}, []string{uuid})
		return
	}

	if strings.TrimSpace(content) == "" {
		log.Println("=== content is empty ===")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " content is empty")}, []string{uuid})
		return
	}

	created_at := time.Now().Format("2006-01-02 15:04:05")

	var message WS_CHAT_MESSAGE_DTO

	client := get_client_by_uuid(clients, uuid)
	message.Content = content
	message.Nickname = client.NICKNAME
	message.Client_number = client.NUMBER
	message.Created_at = created_at

	clients_uuids := get_all_clients_uuids(clients)

	wsSend(WS_CHAT_MESSAGE, message, clients_uuids)

}
