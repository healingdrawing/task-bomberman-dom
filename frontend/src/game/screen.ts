import { screen_prepare } from "./screen_prepare"

/** server controllable */
export enum GameState {
  PLAYER_GAME_OVER = 1,
  START_GAME,
  END_GAME,
}

class GameScreen {

  game_state = GameState.PLAYER_GAME_OVER

  weak_obstacles: Map<string, HTMLDivElement>
  power_ups: Map<string, HTMLDivElement>
  bombs: Map<string, HTMLDivElement>
  players: Map<string, HTMLDivElement>
  explosions: Map<string, HTMLDivElement>

  constructor() {
    this.weak_obstacles = screen_prepare.weak_obstacles
    this.power_ups = screen_prepare.power_ups
    this.bombs = screen_prepare.bombs
    this.players = screen_prepare.players
    this.explosions = screen_prepare.explosions
  }

  /**generates all required objects to display the game and place them on screen*/
  async prepare() {
    screen_prepare.prebuild_game_field()
  }

  game_state_start_game() {
    console.log("=========== game_state_start_game")
    screen_prepare.clear_game_field()
    screen_prepare.build_game_field()
    this.game_state = GameState.START_GAME
  }

  game_state_end_game() {
    this.game_state = GameState.END_GAME
  }

  game_state_player_game_over() {
    this.game_state = GameState.PLAYER_GAME_OVER
  }
}

export const screen = new GameScreen()