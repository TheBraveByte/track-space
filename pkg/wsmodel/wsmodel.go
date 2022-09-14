package wsmodel

import "github.com/yusuf/track-space/pkg/wsconfig"
/*Working with Web Socket*/

type SocketPayLoad struct {
	Condition  string           `json:"condition"`
	Message    string           `json:"message"`
	UserName   string           `json:"username"`
	SocketConn wsconfig.SocketConnection `json:"-"`
}

type SocketResponse struct {
	Condition     string   `json:"condition"`
	Message       string   `json:"message"`
	MessageType   string   `json:"message_type"`
	ConnectedUser []string `json:"connected_user"`
}
