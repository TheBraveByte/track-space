package ws

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yusuf/track-space/pkg/wsconfig"
	"github.com/yusuf/track-space/pkg/wsmodel"
)

// initialize the websocket connection and channel
var webSkChan = make(chan wsmodel.SocketPayLoad)
var Client = make(map[wsconfig.SocketConnection]string)

// UpgradeSocketConn : variable to upgrade ChatRoom controller with a web socket connection
var UpgradeSocketConn = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(rq *http.Request) bool {
		return true
	},
	HandshakeTimeout: 100 * time.Second,
}

// GetDataFromChannel - goroutine method that will get payload data via channel
func GetDataFromChannel() {
	var resp wsmodel.SocketResponse
	for {
		getdata := <-webSkChan
		switch getdata.Condition {
		case "username":
			Client[getdata.SocketConn] = getdata.UserName
			users := GetAllUsers()
			resp.Condition = "username"
			resp.ConnectedUser = users
			BroadCastToAll(resp)
		case "sendMessage":
			resp.Message = fmt.Sprintf("<em>%v</em> : %v", getdata.UserName, getdata.Message)
			resp.Condition = "message"
			BroadCastToAll(resp)

		case "serveroffline":
			resp.Condition = "serveroffline"
			delete(Client, getdata.SocketConn)
			users := GetAllUsers()
			resp.ConnectedUser = users

		}
	}
}

// To get the list of user connected to the web socket
func GetAllUsers() []string {
	var userSlices []string
	for _, user := range Client {
		if user != "" {
			userSlices = append(userSlices, user)
		}
	}
	sort.Strings(userSlices)
	return userSlices
}

// BroadCastToAll - this  functions write a responses to all connected users
func BroadCastToAll(resp wsmodel.SocketResponse) {
	for x := range Client {
		err := x.WriteJSON(resp)
		if err != nil {
			log.Println("error write a json response")
			_ = x.Close()
			delete(Client, x)
			return
		}
	}
}

// SendDataToChannel - goroutine method that will send back response data
// back to the channel
func SendDataToChannel(socketConn *wsconfig.SocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovery from a failed program")
			return
		}
	}()

	var payload wsmodel.SocketPayLoad

	for {
		err := socketConn.ReadJSON(&payload)
		if err != nil {
		} else {
			payload.SocketConn = *socketConn
			webSkChan <- payload
		}

	}
}
