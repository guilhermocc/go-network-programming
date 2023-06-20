## go-network-programming

This repository holds the code samples used during my presentation at Golang SP on June 2023.

### TCP chat
The tcp/chat folder contains a simple chat server and client implementation using TCP sockets.

To run the server:
```bash
$ go run tcp/chat/server/main.go
```

To run the client:
```bash
$ go run tcp/chat/client/main.go
```

The client starts a new connection to the server and waits for user input. When the user types a message and hits enter, the message is sent to the server and broadcasted to all connected clients.

### UDP game

The udp/game folder contains a simple game server and client implementation using UDP sockets. The game is a simple "guess the number" game, where the server picks a random number between 1 and 100 and the client has to guess it.

To run the server:
```bash
$ go run udp/game/server/main.go
```

To run the client:
```bash
$ go run udp/game/client/main.go
```

