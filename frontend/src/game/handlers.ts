import { BroadcastMessage, ChatMessage, ConnectedPlayers, PlayerLifes } from "./types";

class Handlers {
  chat_messages_div: HTMLDivElement;
  constructor() {
    this.chat_messages_div = document.getElementById("chat_view__messages") as HTMLDivElement;
  }

  broadcast_message(data: BroadcastMessage) {
    const message_div = document.createElement("div");
    message_div.classList.add("message_div"); // to style border

    const content_div = document.createElement("div");
    content_div.innerHTML = data.content;
    message_div.appendChild(content_div);

    this.chat_messages_div.insertBefore(message_div, this.chat_messages_div.firstChild);
  }

  chat_message(data: ChatMessage) {
    const message_div = document.createElement("div");
    message_div.classList.add("message_div"); // to style border

    const nickname_div = document.createElement("div");
    nickname_div.innerText = `[${data.nickname}]`;
    nickname_div.classList.add(`color${data.client_number}`);
    message_div.appendChild(nickname_div);

    const time_div = document.createElement("div");
    time_div.innerText = data.created_at;
    message_div.appendChild(time_div);

    const content_div = document.createElement("div");
    content_div.innerText = data.content;
    message_div.appendChild(content_div);
    this.chat_messages_div.insertBefore(message_div, this.chat_messages_div.firstChild);
  }

  connected_players(data: ConnectedPlayers) {
    const connected_players_div = document.getElementById("connected_players") as HTMLDivElement;
    connected_players_div.innerText = `Players: ${data.connected_players}`;
  }

  player_lifes(data: PlayerLifes) {
    const player_lifes_div = document.getElementById("player_lifes") as HTMLDivElement;
    player_lifes_div.innerText = `Lifes: ${data.lifes}`;
  }
}

export const handlers = new Handlers();

