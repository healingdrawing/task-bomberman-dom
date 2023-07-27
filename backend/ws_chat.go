package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WS_GROUP_CHAT_MESSAGE_DTO struct {
	Content    string `json:"content"`
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Created_at string `json:"created_at"`
	Group_id   int    `json:"group_id"`
}

func wsGroupChatMessageHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		return
	}

	_group_id, ok := messageData["group_id"].(float64)
	if !ok {
		log.Println("failed to get group_id from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " failed to get group_id from messageData")}, []string{uuid})
		return
	}
	group_id := int(_group_id)

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

	var message WS_GROUP_CHAT_MESSAGE_DTO

	message.Content = content
	message.Created_at = created_at
	message.Group_id = group_id

	//get all id's of users who is the member of the group of this group chat
	var group_member_ids []int
	log.Println("REFACTOR TO GET ALL UUID'S, BECAUSE THIS IS THE GAME CHAT MESSAGE")
	// get connected user ids

	// get ids of users who is member and connected
	connected_group_member_ids := map[int]int{}
	clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		for _, member_id := range group_member_ids {
			if client.USER_ID == member_id {
				connected_group_member_ids[member_id] = member_id
				break
			}
		}
		return true
	})

	//get all uuids of users who is connected/logged in and member of the group
	var user_uuids []string
	clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		_, ok := connected_group_member_ids[client.USER_ID]
		if ok {
			user_uuids = append(user_uuids, key.(string))
		} else {
			log.Println("user is not connected/logged in anymore")
		}
		return true
	})

	wsSend(WS_GROUP_CHAT_MESSAGE, message, user_uuids)

}
