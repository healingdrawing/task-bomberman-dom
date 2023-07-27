package main

func registerHandlers() {
	// Websocket
	CustomRouter.HandleFunc("/ws", wsConnection)

	// API
	CustomRouter.HandleFunc("/api/user/logout", userLogoutHandler)
	CustomRouter.HandleFunc("/api/user/register", userRegisterHandler)
}
