package main

type WS_RESPONSE_MESSAGE_DTO struct {
	Type WSMT        `json:"type"`
	Data interface{} `json:"data"`
}

type WS_ERROR_RESPONSE_DTO struct {
	Content string `json:"content"`
}

type WS_INFO_RESPONSE_DTO struct {
	Content string `json:"content"`
}

type WS_SUCCESS_RESPONSE_DTO struct {
	Content string `json:"content"`
}
