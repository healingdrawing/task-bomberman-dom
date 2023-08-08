import { BroadcastMessage, ChatMessage, ConnectedPlayers, GameState, MoveVertical, WSMT } from "./types";
import { SendNickname } from "./types";
import { handlers } from "./handlers";
import { screen } from "./screen";

export interface Message {
  type: WSMT;
  data: object;
}

class WebSocketClient {
  private ws: WebSocket | null = null;
  private uuid: string | null = null;
  private readonly serverUrl: string = `ws://localhost:8080/ws`

  private connect(): Promise<void> {
    return new Promise((resolve, reject) => {

      this.ws!.onopen = () => {
        console.log('WebSocket connection established.');
        resolve();
      };

      this.ws!.onerror = (event) => {
        console.error('WebSocket connection failed:', event);
        reject();
      };
    });
  }

  private waitForConnection(): Promise<void> {
    return new Promise((resolve, reject) => {
      const startTime = Date.now();
      const interval = setInterval(() => {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
          clearInterval(interval);
          resolve();
        } else if (Date.now() - startTime > 5000) {
          clearInterval(interval);
          reject(new Error('WebSocket connection took too long to establish.'));
        }
      }, 500);
    });
  }

  private storeUuidInSessionStorage(): void {
    if (this.uuid) {
      sessionStorage.setItem('uuid', this.uuid);
    }
  }

  private handleMessage(message: Message): void {
    switch (message.type) {
      case WSMT.WS_UP:
        console.log('Move up received:', message.data);
        screen.player_move_up(WSMT.WS_UP, message.data as MoveVertical);
        break;
      case WSMT.WS_DOWN:
      case WSMT.WS_LEFT:
      case WSMT.WS_RIGHT:
        // Handle move logic here
        console.log('Move received:', message.data);
        break;
      case WSMT.WS_BOMB:
        // Handle bomb logic here
        console.log('Bomb received:', message.data);
        break;
      case WSMT.WS_EXPLODE:
        // Handle explosion logic here
        console.log('Explode received:', message.data);
        break;
      case WSMT.WS_BROADCAST_MESSAGE:
        // Handle broadcast message logic here
        console.log('Broadcast message received:', message.data);
        handlers.broadcast_message(message.data as BroadcastMessage);
        break;
      case WSMT.WS_CHAT_MESSAGE:
        // Handle chat message logic here
        console.log('Chat message received:', message.data);
        handlers.chat_message(message.data as ChatMessage);
        break;
      case WSMT.WS_CONNECTED_PLAYERS:
        // Handle connected players logic here
        console.log('Connected players received:', message.data);
        handlers.connected_players(message.data as ConnectedPlayers);
        break;
      case WSMT.WS_START_GAME:
        // Handle start game logic here
        console.log('Start game received:', message.data);
        screen.game_state_start_game(message.data as GameState);
        break;
      default:
        console.warn('Unknown message type:', message.type);
    }
  }

  public async initialize(): Promise<void> {
    try {
      this.ws = new WebSocket(this.serverUrl);

      const nickname = document.getElementById("nickname") as HTMLInputElement;

      console.log('1 before this.ws!.onmessage');
      // Server will return uuid after connection is established
      let uuidReceived = false; // Flag to check if UUID has been received
      this.ws!.onmessage = (event) => {
        const message: Message = JSON.parse(event.data);
        console.log('5 before this.uuid message:', message);
        if (!uuidReceived) {
          if (message.type === WSMT.WS_CLIENT_CONNECTED_TO_SERVER) {
            if (
              message.data &&
              typeof message.data === 'object' &&
              'client_uuid' in message.data &&
              'client_number' in message.data
            ) {
              this.uuid = message.data.client_uuid as string;
              this.storeUuidInSessionStorage();
              uuidReceived = true;

              // Send client nickname and uuid to the server
              const initialMessage: Message = {
                type: WSMT.WS_CONNECT_TO_SERVER,
                data: { nickname: nickname.value, client_uuid: this.uuid } as SendNickname,
              };
              this.ws!.send(JSON.stringify(initialMessage));

            } else {
              throw new Error('Failed to extract UUID from the server response.');
            }
          }
        } else {
          // Handle other message types after UUID is received
          this.handleMessage(message);
        }
      };

      console.log('2 Initializing WebSocket client...');
      console.log('3 before this.connect()');
      await this.connect();
      console.log('4 before this.waitForConnection()');
      await this.waitForConnection();
    } catch (error) {
      console.error('Error initializing WebSocket client:', error);
      alert('No free slots!\nPage will be reloaded.\nTry again later..');
      location.reload();
    }
  }

  // uuid will be added automatically to the message.data
  public sendMessage(type: WSMT, data: object): void {
    const extended_data = { ...data, client_uuid: this.uuid };
    const message: Message = {
      type: type,
      data: extended_data,
    };
    console.log('Sending message:', message);
    this.ws!.send(JSON.stringify(message));
  }
}

/** web socket client to interact with server */
export const ws = new WebSocketClient();