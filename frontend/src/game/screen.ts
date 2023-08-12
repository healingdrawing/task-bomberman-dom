import { handlers } from "./handlers"
import { OffsetAnimationParameters, get_from_x, get_from_y, offset_div_animation } from "./screen_offset_animation"
import { screen_prepare } from "./screen_prepare"
import { BombXY, EndGame, ExplodeBomb, GameState, HidePowerUp, MoveDx, MoveDy, PlayerLifes, WSMT } from "./types"

/** server controllable */
export enum GameStateValue {
  PLAYER_GAME_OVER = 1,
  START_GAME,
  END_GAME,
}

export class GameScreen {

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

  game_state_end_game(data: EndGame, uuid: string) {
    this.game_state_value = GameStateValue.END_GAME
    if (data.winner_uuid === uuid) {
      alert("!!!Victory!!!")
    }
    else {
      alert("Congratulations, you failed again. Amazing stability!")
    }
    location.reload()
  }

  game_state_player_game_over() {
    this.game_state_value = GameStateValue.PLAYER_GAME_OVER
  }

  /** at the moment only checks player game over */
  check_game_state_value(data: PlayerLifes) {
    if (data.lifes < 1) {
      console.log("=========== game_state_player_game_over fired")
      this.game_state_player_game_over()
    }
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

    // remove weak obstacles and show power ups
    explode_data.destroy_xy.forEach((xy, index) => {
      const weak_obstacle = this.weak_obstacles.get(xy)
      if (weak_obstacle) {
        weak_obstacle.classList.remove(`weak_obstacle`)
        weak_obstacle.classList.add(`none`)

        if (explode_data.power_up_effect[index] !== "0") {
          const power_up = this.power_ups.get(xy)
          if (power_up) {
            power_up.classList.remove(`none`)
            power_up.classList.add(`power_up${explode_data.power_up_effect[index]}`)
          }
        }

      }
    })

    // execute explosion
    const last = explode_data.cells_xy.length - 1
    explode_data.cells_xy.forEach((cell, index) => {
      const explosion = this.explosions.get(cell)
      if (explosion) {
        explosion.classList.remove(`none`)
        explosion.classList.remove(`explosion`)

        //todo: very hungry for performance
        void explosion.offsetWidth; // must be on every item(checked) to overwrite animation. Affect performance, but after some time, stops affect performance
        /*
        but otherwise overcrossed animations, must be managed by removing
        of the classes with delay, which can damage new animation raised
        before first one completed, and only garbaging of the screen
        by z-index stacking of the divs can be the solution,
        with again removing with delay. And task requires as possible less layers, do not know what is worse, pile of shit on screen, or some decreasing of fps but with simple code. Let is hope it will pass audits.
        */

        explosion.classList.add(`explosion`)
      }
    })
  }

  hide_power_up(hide_data: HidePowerUp) {
    const power_up = this.power_ups.get(hide_data.cell_xy)
    if (power_up) {
      power_up.classList.remove(`power_up${hide_data.effect}`)
      power_up.classList.add("none")
    }
  }

}

export const screen = new GameScreen()