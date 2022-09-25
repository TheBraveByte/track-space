package wsconfig

import "github.com/gorilla/websocket"

// SocketConnection : WebSocket connection setup on the server side with Gorilla
type SocketConnection struct {
	*websocket.Conn
}
