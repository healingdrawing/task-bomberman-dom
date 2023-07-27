/**web socket message type. The same values as server side */
export enum WSMT {
  WS_ERROR_RESPONSE = "error_response",
  WS_CHAT_MESSAGE = "chat_message",

  WS_CONNECT_TO_SERVER = "connect_to_server",
  WS_KEEP_CONNECTION = "keep_connection",
  WS_KILL_CONNECTION = "kill_connection",
  WS_CLIENT_CONNECTED_TO_SERVER = "client_connected_to_server",
  WS_STILL_CONNECTED = "still_connected",
  WS_BROADCAST_MESSAGE = "broadcast_message",
  WS_START_GAME = "start_game",
  WS_END_GAME = "end_game",
  WS_PLAYER_GAME_OVER = "player_game_over",
  WS_UP = "up",
  WS_DOWN = "down",
  WS_LEFT = "left",
  WS_RIGHT = "right",
  WS_STAND = "stand",
  WS_BOMB = "bomb",
  WS_EXPLODE = "explode",
  WS_HIDE_POWER_UP = "hide_power_up"
}