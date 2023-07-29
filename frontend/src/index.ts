import { router, store, events } from './framework/framework';
import { WebSocketClient } from './game/ws';

const websockets: (WebSocket | null)[] = [];
const serverUrl = `ws://localhost:8080/ws`;
var client: WebSocketClient; //todo: later maybe hide this from global scope

async function connect_to_game(event: Event) {
  event.preventDefault(); // Prevent the default form submission

  const inputField = document.getElementById("nickname") as HTMLInputElement;

  if (!inputField.checkValidity()) {
    // Display an error message or take appropriate action
    alert(`🤦 Before play this game you should learn how to properly type on keyboard.
    Please leave this place, it can be too danger for your virgin brain,
    which is not able to choose nickname correctly.`);
    return;
  } else {
    client = new WebSocketClient(serverUrl, inputField.value);
    await client.initialize();
  }

  const inputValue = inputField.value;

  // Process the input value
  alert("Input value: " + inputValue);
}

/**initialization */
(() => {
  const connect_to_game_button = document.getElementById("connect_to_game") as HTMLButtonElement;
  events.on("click", connect_to_game_button, connect_to_game); // usage of mini-framework 🙂 mission complete

})()