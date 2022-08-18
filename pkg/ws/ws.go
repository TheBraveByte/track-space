package ws

import (
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yusuf/track-space/pkg/model"
)

// initialize the websocket connection and channel
var WebSkChan = make(chan model.SocketPayLoad)

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
	var resp model.SocketResponse
	for {
		getdata := <- WebSkChan
		switch getdata.Condition{
		case "username":
			model.Client[getdata.SocketConn] = getdata.UserName
			users := GetAllUsers()
			resp.Condition= "username"
			resp.ConnectedUSer= users 
		}
	}
}
//To get the list of user connected to the web socket
func GetAllUsers() []string{
	var userSlices []string
	for _, user :=  range model.Client{
		if user !=  ""{
			userSlices = append(userSlices, user)
		}

	}
	sort.Strings(userSlices)
	return userSlices
}

//BroadCastToAll - this  functions write a responses to all connected users
func BroadCastToAll(resp model.SocketResponse) {
	for x :=range model.Client{
		err := x.WriteJSON(resp)
		if err != nil{
			log.Println("error write a json response")
			_ = x.Close()
			return
		}
	}
}
// SendDataToChannel - goroutine method that will send back response data
// back to the channel
func SendDataToChannel(socketConn *model.SocketConnection, c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovery from a failed program")
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{})
			return
		}
	}()

	var payload model.SocketPayLoad

	for {
		err := socketConn.ReadJSON(&payload)
		if err != nil {
			// _ = c.AbortWithError(http.StatusInternalServerError, gin.Error{})
		} else {
			payload.SocketConn = *socketConn
			model.WebSkChan <- payload
		}

	}
}