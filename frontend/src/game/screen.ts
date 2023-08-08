import { handlers } from "./handlers"
import { screen_prepare } from "./screen_prepare"
import { GameState, MoveVertical, WSMT } from "./types"

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
    handlers.player_lifes({ lifes: 3 })
    this.game_state_value = GameStateValue.START_GAME
  }

  game_state_end_game() {
    this.game_state_value = GameStateValue.END_GAME
  }

  game_state_player_game_over() {
    this.game_state_value = GameStateValue.PLAYER_GAME_OVER
  }



  /* player animation section */
  player_move_up(move_type: WSMT, move_data: MoveVertical) {
    console.log("inside player_move", move_type)
    const player = this.players.get(move_data.number.toString())
    if (player) {
      const params = {
        // fromX: 0,
        // toX: 0,
        // fromY: 5 * this.cspx,
        // toY: (move_data.target_y - 1) * this.cspx,
        // speed: screen.cspx * (move_data.turbo ? 2 : 1),
        fromX: 0,
        toX: 0,
        fromY: 500,
        toY: 0,
        speed: 100,
      } as OffsetAnimationParameters
      offset_div_animation(player, params.fromX, params.toX, params.fromY, params.toY, 1000)
    }

  }
}

function offset_div_animation(
  div: HTMLElement,
  fromX: number,
  toX: number,
  fromY: number,
  toY: number,
  duration: number
) {
  const intervalDuration = 16; // Approximate interval duration for 60 FPS
  const frames = Math.ceil(duration / intervalDuration);

  const xStep = (toX - fromX) / frames;
  const yStep = (toY - fromY) / frames;

  let frameCount = 0;

  const animationInterval = setInterval(() => {
    frameCount++;
    if (frameCount >= frames) {
      clearInterval(animationInterval);
    }

    const currentX = fromX + xStep * frameCount;
    const currentY = fromY + yStep * frameCount;

    div.style.left = `${currentX}px`;
    div.style.top = `${currentY}px`;
  }, intervalDuration);

  setTimeout(() => {
    clearInterval(animationInterval);
  }, duration);
}

interface OffsetAnimationParameters {
  fromX: number;
  toX: number;
  fromY: number;
  toY: number;
  speed: number;
}

export const screen = new GameScreen()