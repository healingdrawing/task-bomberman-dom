import { router, store, events } from './framework/framework';
import { State } from './framework/store';
import { prebuild_game_field } from './game/screen';
import { SendChatMessage, WSMT } from './game/types';
import { WebSocketClient } from './game/ws';

// switch views on screen, when connected to server
store.setState({
  first_screen_visible: true,
  second_screen_visible: false,
});
store.subscribe((state: State) => {
  const firstScreenDiv = document.getElementById("first_screen") as HTMLDivElement;
  const secondScreenDiv = document.getElementById("second_screen") as HTMLDivElement;
  // Update visibility based on state
  firstScreenDiv.style.display = state.first_screen_visible ? "block" : "none";
  secondScreenDiv.style.display = state.second_screen_visible ? "block" : "none";
});

const serverUrl = `ws://localhost:8080/ws`
var client: WebSocketClient //todo: later maybe hide this from global scope

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
    const connect_to_game_button = document.getElementById("connect_to_game") as HTMLButtonElement;
    events.off("click", connect_to_game_button, connect_to_game, "connect_to_game_button");

    client = new WebSocketClient(serverUrl, inputField.value);
    await client.initialize();

    store.setState({
      first_screen_visible: !store.getState().first_screen_visible,
      second_screen_visible: !store.getState().second_screen_visible,
    });

  }
}

function chat_message() {
  const inputField = document.getElementById("chat_view__input__message") as HTMLInputElement;
  const inputValue = inputField.value;
  const message = {
    content: inputValue
  } as SendChatMessage
  inputField.value = ""
  client.sendMessage(WSMT.WS_CHAT_MESSAGE, message)
}

/**initialization */
(async () => {
  console.log("prepare environment")
  await prebuild_game_field();
  const connect_to_game_button = document.getElementById("connect_to_game") as HTMLButtonElement;
  events.on("click", connect_to_game_button, connect_to_game); // usage of mini-framework 🙂 mission complete
  const chat_message_button = document.getElementById("chat_view__input__send") as HTMLButtonElement;
  events.on("click", chat_message_button, chat_message);
})()