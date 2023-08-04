import { WSMT } from "./types";
import { GameState, screen } from "./screen";
import { ws } from "./ws";

export class Controller {
  // Create a Set to keep track of pressed keys
  private pressedKeys = new Set<string>();

  private listenedKeys = new Map<string, [WSMT, WSMT]>(
    [
      ["Enter", [WSMT.WS_BOMB_OFF, WSMT.WS_BOMB_ON]],
      ["ArrowUp", [WSMT.WS_UP_OFF, WSMT.WS_UP_ON]],
      ["ArrowDown", [WSMT.WS_DOWN_OFF, WSMT.WS_DOWN_ON]],
      ["ArrowLeft", [WSMT.WS_LEFT_OFF, WSMT.WS_LEFT_ON]],
      ["ArrowRight", [WSMT.WS_RIGHT_OFF, WSMT.WS_RIGHT_ON]],
    ]
  )

  // Function to handle keydown event
  handleKeyDown = (event: KeyboardEvent) => {
    const key = event.key;
    // Check if the key is not pressed, and the key is one of the listened keys
    if (screen.game_state == GameState.START_GAME && !this.pressedKeys.has(key) && this.listenedKeys.has(key)) {
      ws.sendMessage(this.listenedKeys.get(key)![1], {})
      // Add the key to the Set
      this.pressedKeys.add(key);
      console.log("key down", key) // todo: remove
    }
  }

  // Function to handle keyup event
  handleKeyUp = (event: KeyboardEvent) => {
    const key = event.key;

    // Check if the key is already pressed, and the key is one of the listened keys
    if (screen.game_state == GameState.START_GAME && this.listenedKeys.has(key)) {
      ws.sendMessage(this.listenedKeys.get(key)![0], {})
      // Remove the key from the Set
      this.pressedKeys.delete(key);
      console.log("key up", key) // todo: remove
    }
  }
}

// Create an instance of the Controller class
/**keyboard keys press managing */
export const controller = new Controller();
