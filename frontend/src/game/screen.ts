import { screen_prepare } from "./screen_prepare"

class GameScreen {
  /**generates all required objects to display the game and place them on screen*/
  async prepare() {
    await screen_prepare.prebuild_game_field()
  }
}

export const screen = new GameScreen()