package wsconfig

import "github.com/gorilla/websocket"

type SocketConnection struct {
	*websocket.Conn
}
