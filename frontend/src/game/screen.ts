import { handlers } from "./handlers"
import { OffsetAnimationParameters, get_from_x, get_from_y, offset_div_animation } from "./screen_offset_animation"
import { screen_prepare } from "./screen_prepare"
import { BombXY, ExplodeBomb, GameState, MoveDx, MoveDy, WSMT } from "./types"

/** server controllable */
export enum GameStateValue {
  PLAYER_GAME_OVER = 1,
  START_GAME,
  END_GAME,
}

class GameScreen {

  game_state_value = GameStateValue.PLAYER_GAME_OVER

  /** cell size in pixels, 96px inside ScreenPrepare */
  cspx: number

  weak_obstacles: Map<string, HTMLDivElement>
  power_ups: Map<string, HTMLDivElement>
  bombs: Map<string, HTMLDivElement>
  players: Map<string, HTMLDivElement>
  explosions: Map<string, HTMLDivElement>

  constructor() {
    this.cspx = screen_prepare.cspx
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

  game_state_start_game(state: GameState) {
    console.log("=========== game_state_start_game")
    screen_prepare.clear_game_field()
    screen_prepare.build_game_field(state)
    handlers.player_lifes({ lifes: 3 }) //todo: hardcoded, to improve performance
    this.game_state_value = GameStateValue.START_GAME
  }

  game_state_end_game() {
    this.game_state_value = GameStateValue.END_GAME
  }

  game_state_player_game_over() {
    this.game_state_value = GameStateValue.PLAYER_GAME_OVER
  }



  /* player animation/offset section */
  /** move player along vertical screen direction */
  player_move_dy(move_type: WSMT, move_data: MoveDy) {
    console.log("inside player_move", move_type)
    const player = this.players.get(move_data.number.toString())
    if (player) {
      const from_x = get_from_x(player, screen.cspx)
      const from_y = get_from_y(player, screen.cspx)
      const params = {
        div: player,
        player_number: move_data.number,
        fromX: from_x,
        toX: from_x,
        fromY: from_y,
        toY: move_data.target_y,
        cells_per_second: move_data.turbo ? 2 : 1,
        cspx: screen.cspx,
      } as OffsetAnimationParameters
      offset_div_animation(params)
    }
  }

  player_move_dx(move_type: WSMT, move_data: MoveDx) {
    console.log("inside player_move", move_type)
    const player = this.players.get(move_data.number.toString())
    if (player) {
      const from_x = get_from_x(player, screen.cspx)
      const from_y = get_from_y(player, screen.cspx)
      const params = {
        div: player,
        player_number: move_data.number,
        fromX: from_x,
        toX: move_data.target_x,
        fromY: from_y,
        toY: from_y,
        cells_per_second: move_data.turbo ? 2 : 1,
        cspx: screen.cspx,
      } as OffsetAnimationParameters
      offset_div_animation(params)
    }
  }

  player_bomb_xy(bomb_data: BombXY) {
    console.log("inside player bomb xy")
    const bomb = this.bombs.get(bomb_data.target_xy)
    if (bomb) {
      bomb.classList.remove(`none`)
      bomb.classList.add(`player${bomb_data.number}_bomb`)
    }
  }

  explode_bomb(explode_data: ExplodeBomb) {
    console.log("inside explode bomb")

    // remove bomb

    const bomb = this.bombs.get(explode_data.cells_xy[0])
    if (bomb) {
      bomb.classList.remove(`player${explode_data.number}_bomb`)
      bomb.classList.add(`none`)
    }

    // execute explosion
    const last = explode_data.cells_xy.length - 1
    explode_data.cells_xy.forEach((cell, index) => {
      const explosion = this.explosions.get(cell)
      if (explosion) {
        explosion.classList.remove(`none`)
        explosion.classList.remove(`explosion`)

        //todo: very hungry for performance
        void explosion.offsetWidth; // must be on every item to overwrite animation. Affect performance

        explosion.classList.add(`explosion`)
      }
    })
  }

}

export const screen = new GameScreen()