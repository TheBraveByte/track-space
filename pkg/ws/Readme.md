# WS

Package **ws** provides a basic implementation of a websocket server. It uses the Gorilla WebSocket library to handle the websocket communication. The package is used to initialize the websocket connection and channel, upgrade the ChatRoom controller with a web socket connection, get data via a channel, and send back response data to the channel.

### Features

* **UpgradeSocketConn**: variable to upgrade ChatRoom controller with a web socket connection.

* **GetDataFromChannel()**: goroutine method that will get payload data via channel.

* **GetAllUsers()**: function that returns a list of users connected to the websocket.

* **BroadCastToAll(resp wsmodel.SocketResponse)**: function that writes a response to all connected users.

* **SendDataToChannel(socketConn *wsconfig.SocketConnection)**: goroutine method that will send back response data back to the channel.

### Usage
To use the package, import it in your application:

```go
import (
	"github.com/yusuf/track-space/pkg/ws"
)

// To upgrade the ChatRoom controller with a web socket connection, use the UpgradeSocketConn variable:


http.HandleFunc("/chat-room", func(w http.ResponseWriter, r *http.Request) {
    conn, err := ws.UpgradeSocketConn.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    // Handle the websocket connection
    go ws.SendDataToChannel(wsconfig.NewSocketConnection(conn))
    go ws.GetDataFromChannel()
})

```

In the above example, `SendDataToChannel` and `GetDataFromChannel` are used to handle the websocket communication.

### Contributions

Contributions to the package are welcome. You can submit an issue to report a bug or suggest an enhancement. You can also fork the package and submit a pull request with your changes.
