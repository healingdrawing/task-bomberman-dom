/**web socket message types. The same values as server(backend) side */
export enum WSMT {
  WS_ERROR_RESPONSE = "error_response",
  WS_CHAT_MESSAGE = "chat_message",

  WS_CLIENT_CONNECTED_TO_SERVER = "client_connected_to_server", // first response with uuid
  WS_CONNECT_TO_SERVER = "connect_to_server", // first message to server with nickname
  WS_KEEP_CONNECTION = "keep_connection", // keep connection alive
  WS_KILL_CONNECTION = "kill_connection", // todo: not sure it needed
  WS_STILL_CONNECTED = "still_connected", // response from server, that client is still connected
  WS_BROADCAST_MESSAGE = "broadcast_message",
  WS_START_GAME = "start_game",
  WS_END_GAME = "end_game",
  WS_PLAYER_GAME_OVER = "player_game_over",

  // incoming from server to control items on screen
  WS_UP = "up",
  WS_DOWN = "down",
  WS_LEFT = "left",
  WS_RIGHT = "right",
  WS_STAND = "stand",
  WS_BOMB = "bomb",
  WS_EXPLODE = "explode",
  WS_HIDE_WEAK_OBSTACLES = "hide_weak_obstacles", // todo: also must includes show power ups
  WS_HIDE_POWER_UP = "hide_power_up",

  // outgoing from client to server to send pressed/released keys
  WS_UP_ON = "up_on", // arrow up pressed
  WS_UP_OFF = "up_off", // arrow up released
  WS_DOWN_ON = "down_on",
  WS_DOWN_OFF = "down_off",
  WS_LEFT_ON = "left_on",
  WS_LEFT_OFF = "left_off",
  WS_RIGHT_ON = "right_on",
  WS_RIGHT_OFF = "right_off",
  WS_BOMB_ON = "bomb_on", // Enter
  WS_BOMB_OFF = "bomb_off",

  WS_PLAYER_LIFES = "player_lifes",
  WS_CONNECTED_PLAYERS = "connected_players", // number of connected players
}

// here uuid added manually, because this is straight message to server, without use "sendMessage" function. Otherwise uuid will be added in sendMessage function automatically
export interface SendNickname {
  nickname: string
  client_uuid: string
}

export interface BroadcastMessage {
  content: string
  client_number: number
}

export interface SendChatMessage {
  content: string
}

export interface ChatMessage {
  content: string
  nickname: string
  client_number: number
  created_at: string
}

export interface ConnectedPlayers {
  connected_players: string
}

/** server sends number:number field also, which is player number(not used at the moment) */
export interface PlayerLifes {
  lifes: number
}

/* game state */

export interface Player {
  number: number
  x: number
  y: number
  target_x: number
  target_y: number
  turbo: boolean
  lifes: number
}

export interface Bomb {
  player_number: number
  show: boolean
}

export interface PowerUp {
  effect: string
  show: boolean
}

/**  do not convert to maps etc structures, because no short way to do this, only manuall convertion, which is awful. Also iteration through the object can be used*/
export interface GameState {
  // players: object
  players: { [key: string]: Player };
  bombs: object
  explosions: object
  power_ups: object
  weak_obstacles: object
}

export interface StartGame {
  state: GameState
}

export interface MoveDy {
  number: number
  target_y: number
  turbo: boolean
}

export interface MoveDx {
  number: number
  target_x: number
  turbo: boolean
}

export interface BombXY {
  number: number
  target_xy: string
}

export interface ExplodeBomb {
  number: number
  cells_xy: string[]
  destroy_xy: string[] // destroy xy used for power_up too and effect is "1" ... "3" for animations
  power_up_effect: string[]
}

export interface HidePowerUp {
  cell_xy: string
  effect: string
}