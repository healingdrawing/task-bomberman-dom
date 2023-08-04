package main

func registerHandlers() {
	// Websocket
	CustomRouter.HandleFunc("/ws", wsConnection)
}
