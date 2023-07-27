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
	WS_ERROR_RESPONSE   WSMT = "error_response"
	WS_INFO_RESPONSE    WSMT = "info_response"
	WS_SUCCESS_RESPONSE WSMT = "success_response"

	WS_COMMENT_SUBMIT WSMT = "comment_submit"
	WS_COMMENTS_LIST  WSMT = "comments_list"

	WS_PRIVATE_CHAT_USERS_LIST WSMT = "private_chat_users_list"
	WS_PRIVATE_CHAT_MESSAGE    WSMT = "private_chat_message"
	WS_GROUP_CHAT_MESSAGE      WSMT = "group_chat_message"

	WS_FOLLOW_REQUEST_RESPONSE WSMT = "follow_request_response" // NEW when broadcast the follow request from fan user, to idol user. Not implemented yet
	WS_FOLLOW_REQUEST_REJECT   WSMT = "follow_request_reject"
	WS_FOLLOW_REQUEST_ACCEPT   WSMT = "follow_request_accept"
	WS_FOLLOW_REQUESTS_LIST    WSMT = "follow_requests_list"

	WS_POST_SUBMIT WSMT = "post_submit"
	WS_POSTS_LIST  WSMT = "posts_list"

	WS_GROUPS_LIST     WSMT = "groups_list"     // FOR USER MEMBERSHIP
	WS_GROUPS_ALL_LIST WSMT = "groups_all_list" // TO DISCOVER ALL GROUPS

	WS_GROUP_SUBMIT              WSMT = "group_submit"
	WS_GROUP_POST_SUBMIT         WSMT = "group_post_submit"
	WS_GROUP_POSTS_LIST          WSMT = "group_posts_list"
	WS_GROUP_POST_COMMENT_SUBMIT WSMT = "group_post_comment_submit"
	WS_GROUP_POST_COMMENTS_LIST  WSMT = "group_post_comments_list"
	WS_GROUP_JOIN                WSMT = "group_join"

	WS_GROUP_INVITES_SUBMIT WSMT = "group_invites_submit"
	WS_GROUP_INVITE_ACCEPT  WSMT = "group_invite_accept"
	WS_GROUP_INVITE_REJECT  WSMT = "group_invite_reject"
	WS_GROUP_INVITES_LIST   WSMT = "group_invites_list" // FOR USER, many groups

	WS_GROUP_REQUEST_SUBMIT WSMT = "group_request_submit"
	WS_GROUP_REQUEST_ACCEPT WSMT = "group_request_accept"
	WS_GROUP_REQUEST_REJECT WSMT = "group_request_reject"
	WS_GROUP_REQUESTS_LIST  WSMT = "group_requests_list"

	WS_GROUP_EVENT_SUBMIT            WSMT = "group_event_submit"
	WS_GROUP_EVENTS_LIST             WSMT = "group_events_list"
	WS_GROUP_EVENT_PARTICIPANTS_LIST WSMT = "group_event_participants_list" // todo: probably artefact
	WS_GROUP_EVENT_GOING             WSMT = "group_event_going"
	WS_GROUP_EVENT_NOT_GOING         WSMT = "group_event_not_going"

	WS_USER_CHECK                    WSMT = "user_check"
	WS_USER_FOLLOWING_LIST           WSMT = "user_following_list"
	WS_USER_FOLLOWERS_LIST           WSMT = "user_followers_list"
	WS_USER_FOLLOW                   WSMT = "user_follow"
	WS_USER_LOGIN                    WSMT = "user_login"
	WS_USER_LOGOUT                   WSMT = "user_logout"
	WS_USER_POSTS_LIST               WSMT = "user_posts_list"
	WS_USER_GROUP_POSTS_LIST         WSMT = "user_group_posts_list"
	WS_USER_GROUPS_FRESH_EVENTS_LIST WSMT = "user_groups_fresh_events_list"
	WS_USER_PRIVACY                  WSMT = "user_privacy"
	WS_USER_PROFILE                  WSMT = "user_profile"
	WS_USER_REGISTER                 WSMT = "user_register"
	WS_USER_UNFOLLOW                 WSMT = "user_unfollow"

	WS_USER_VISITOR_STATUS       WSMT = "user_visitor_status"
	WS_USER_GROUP_VISITOR_STATUS WSMT = "user_group_visitor_status"
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
