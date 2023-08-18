# gameserver-websocket

This repos implements a game server in Golang.
It uses websockets and a simple in-memory pub/sub broker to share messages with clients.

## Usage

### Install

```sh
make download
```

### Run server

Run the game server with this command:

```sh
make server
```

### Run test client

Run a test game client with this command:

```sh
make client
```

## Client Server protocol

### Websocket

The game server runs a pubsub broker over a websocket server.
The default port is 8000.

The clients connect to the server and send RPC commands.
Each RPC command is made with a `method` and a payload.

### Publish

Publications are messages sent by the game server to connected clients (eg players)

#### Topics

* General game topic
This topic is used for actions related to players actions (player registrations, create or join game, ...)

* Dedicated game topics

These topics are used for actions related to a single game.
Each time a game is created, a unique and dedicated topics is created.

#### Message format

Messages are encoded in JSON with the following schema:
```json
{
    "type": string,
    "actor": string,
    "id": string,
    "data" string,
}
```

The following types are:
* `information`: general information related to the server or a game
* `creation`: new game in the server
* `start`: game started for example
* `join`: player (name provided in data) joined a game
* `registration`: new player registered (name provided in data)

### RPC

#### Methods

* `register`: handles new player registration
* `unregister`: removes a player from registry
* `listAll`: returns the list of all players
...


## Game server Actions

Here's the sequence of actions that includes the ability for the game server to host multiple games simultaneously, allowing players to create new games, join existing games, and record player scores for statistical purposes:

### 1. Server Setup

* The game server initializes and prepares the game environment.
* It sets up the game rules, board, and any necessary configurations.
* The server waits for incoming connections from clients.

### 2. Client Connection

* A client connects to the game server.
* The server verifies the client's credentials (if applicable) and establishes a connection.

### 3. Player Registration

* The client sends a registration request to the server.
* The server validates the registration request and assigns a unique identifier (e.g., player ID) to the client.
* The server m* ay also assign initial game resources or perform any other necessary initialization for the player.

### 4. Create New Game

* A player who wants to create a new game sends a request to the server.
* The server creates a new game instance and assigns a unique game ID to it.
* The player becomes the host of the new game and is automatically joined to it.

### 5. Join Existing Game

* A player who wants to join an existing game sends a request to the server, specifying the game ID they wish to join.
* The server verifies the game ID and checks if there are available slots for new players in the requested game.
* If there are available slots, the server adds the player to the requested game.
* If there are no available slots or the game ID is invalid, the server sends an error message to the player.

### 6. Start Game

* Once all desired players have joined a game, the host sends a start game request to the server.
* The server signals the start of the game and sends a game start message to all connected clients participating in that game.

### 7. Turn-based Gameplay Loop

The server determines which player's turn it is and broadcasts this information to all clients in the game.
The current player receives a turn notification from the server.

### 8. Player Action

The current player's client sends an action request to the server, specifying the intended action to take during their turn.
* The server validates the action and checks if it complies with the game rules.
* If the action is valid, the server updates the game state accordingly and broadcasts the updated state to all clients in the game.
* If the action is invalid, the server sends an error message to the client and prompts them to retry.

### 9. Game State Update

* After each player's action, the server updates the game state and broadcasts the updated state to all clients in the game.
* The clients receive the updated game state and render it on their respective interfaces.

### 10. Next Turn

* The server determines the next player's turn based on the game rules or any other logic.
* The server broadcasts the turn information to all clients in the game, indicating the next player's turn.

### 11. Repeat Steps 8-10

* Steps 8-10 are repeated for each subsequent turn until a game-ending condition is met (e.g., victory, draw, or a specified number of turns).

### 12. End Game

* Once the game-ending condition is met, the server updates the game state to reflect the final result.
* The server records the players' scores for statistical purposes.
* The server sends an end game message to all clients in the game, indicating the result (e.g., winner, loser, draw) and displaying the final scores.
* The clients receive the end game message and display the final result and scores on their interfaces.

### 13. Cleanup and Statistics

* The server performs any necessary cleanup tasks for the ended game.
* The server updates and maintains the player statistics (scores, win/loss ratio, etc.) for future reference.

### Optional: New Game Request

If desired, the clients can send a new game request to the server, indicating their intention to start another game.
Steps 4-13 are repeated to create and play a new game.
With this updated sequence of actions, the game server can host multiple games simultaneously, allowing players to create new games, join existing games, and record player scores for statistical purposes at the end of each game.

## Todo

* [ ] Configuration for port and allowed origin
* [ ] plugin for RPC management