# WS

WS is a package that import WSconfig package which provides a WebSocket connection setup on the server side with Gorilla. It allows you to easily create and manage WebSocket connections for your applications.

### Installation
To install Wesocket package, you can use the go get command:

go get github.com/gorilla/websocket

### Usage
To use WSConfig, you first need to import it into your Go program:
go get github.com/yusuf/pkg/wsconfig

```go
conn, err := wsconfig.New(url)
if err != nil {
    log.Fatal(err)
}
defer conn.Close()
The url argument is a string that specifies the URL of the WebSocket server.

Once you have established a connection, you can send and receive messages using the WriteMessage and ReadMessage functions:

lua
Copy code
err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, world!"))
if err != nil {
    log.Fatal(err)
}

_, message, err := conn.ReadMessage()
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(message))
```

### Features

#### WSConfig provides the following features:

#### WebSocket connection setup
WSConfig makes it easy to set up WebSocket connections on the server side with Gorilla. You can create a new connection by calling the wsconfig.New function.

#### Message sending and receiving
Once you have established a connection, you can send and receive messages using the WriteMessage and ReadMessage functions.

#### Connection closing
WSConfig provides a Close function that can be used to close the WebSocket connection when you are finished with it. This function should be called in a defer statement to ensure that the connection is always properly closed.

Examples
Here's an example of using WSConfig to create a WebSocket connection:

```go

package main

import (
    "fmt"
    "log"

    "github.com/username/wsconfig"
    "github.com/gorilla/websocket"
)

func main() {
    url := "ws://localhost:8080/ws"
    conn, err := wsconfig.New(url)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, world!"))
    if err != nil {
        log.Fatal(err)
    }

    _, message, err := conn.ReadMessage()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(message))
}
```
#### Contributing
If you find a bug or would like to contribute to WSConfig, please submit an issue or a pull request on the GitHub repository: https://github.com/username/wsconfig.
