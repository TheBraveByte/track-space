package wsmodel

import "github.com/yusuf/track-space/pkg/wsconfig"

/*Working with Web Socket*/

/*
SocketPayLoad : this model helps to handle input or content sent from the
Client-side of the application during RTC(real time communication) using
Web socket ( Reconnecting web socket ).
*/
type SocketPayLoad struct {
	Condition  string                    `json:"condition"`
	Message    string                    `json:"message"`
	UserName   string                    `json:"username"`
	SocketConn wsconfig.SocketConnection `json:"-"`
}

/*
SocketResponse : this model helps to handle output response sent from the
Server-side of the application based on the payload received from the client side.
*/
type SocketResponse struct {
	Condition     string   `json:"condition"`
	Message       string   `json:"message"`
	MessageType   string   `json:"message_type"`
	ConnectedUser []string `json:"connected_user"`
}
