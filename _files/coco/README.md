MESSAGE TYPES TO INTERACT BETWEEN CLIENTS AND SERVER.
All messages are WS based instead of marked HTTP

The first "message" is "join_server". Only this "message" is HTTP based request, other messages is websocket based, after this first http request, this http request will be converted to websocket connection(like wire/channel/tunnel etc). The uuid goten after first message must be included in every message from client, to allow server know who is sender

= =  
= CLIENT sends request messages to server: =  
= =  

- `"connect_to_server"` - the HTTP(not a WS websocket) based  request to server. Must be sent from the first game ui screen, where member name is necessarily to recognize players/members.

- `"keep_connection"` - must be sent every 0.5 second, otherwise server will kill the ws connection

- `"kill_connection"` - must be sent, if there was no "still_connected" server response more than one second, or if user press the button "quit game". After that first ui screen must be displayed

- `"chat_message"` - send to chat

- `"player_game_over"` - it must be button "R.I.P". After press it player dies, the "r.i.p" stone on his place. Player becomes an observer.

= =  
= SERVER sends response messages to clients:  
= =

- `"client_connected_to_server"` - inlcudes uuid created on server side to identify user. uuid must be saved in client side and used inside every message later. After that message arrives, the next screen appears. Then user can use chat and try to join the game.

- `"still_connected"` - must be sent, as response to client message "keep_connection", as confirmation to client that connection is still alive. Otherwise client will sent the "kill_connection" message , and returns to first ui view before connection was established. The "kill_connection" message will be sent if client does not have "still_connected" message response, more then one second.

- `"broadcast_message"` - response to every client "chat_message". Also used for notifications: "server_notification", "new_user_connected", "user_disconnected", "users_number", "wait_seconds", "prepare_seconds".
The "wait_seconds" (20 seconds for players) and "prepare_seconds" (10 seconds before game starts) will be sent every second, like countdown.

- `"start_game"` - command for clients to show game screen, and start the game play

- `"end_game"` - command for clients end the gameplay and show the waiting screen, happens after winner is determined(10 seconds later). Before this, every user can continue to see the game play process, even after loose x3 lifes. Like observer.

- `"player_game_over"` - no lives anymore, or player pressed button "R.I.P". Present player position nearest cell center must be replaced by "R.I.P" stone

= =  
= gameplay control messages. The names/types are the same for CLIENT and SERVER  
= so client send request with this type, and server responds with the same type  
= if action is not possible, server does not respond, f.e. try move on the wall  
= these messages responses were include the number of player etc, to affect  
= =  

- `"up"` or `"down"` or `"left"` or `"right"` - to try move player along level map. SERVER responds x10 times per second, until move completed, by new coordinates of player destination (x,y). To make move animation more SMOOTH/POLISHED, without jumps, perhaps on the client side will be implemented additional animation method for player move. It will starts the new smooth move, for every new server respond for move player, from the present player position to the destination coordinates in server message. It will happens for each player on screen.

- `"stand"` - sends by client every time when `"up"`, `"down"`, `"left"`, `"right"` UNPRESSED, to say to server do not continue moving, after previous move completed

- `"bomb"` - to place the bomb in the map nearest cell where the player stands, if the cell is empty(no bomb placed before, ... + check not a brick, not a wall(perhaps impossible cases))

- `"explode"` - CLIENT send it, to explode the bomb of the player(let it be one 🙂 ). SERVER responds with instructions which bomb must be replaced by explosion animation. Which bricks must disappear after explosion. Which break must be replaced by power up(and the type of power up of course).

- `"hide_power_up"` - SERVER response to CLIENTS when power up was taken by user