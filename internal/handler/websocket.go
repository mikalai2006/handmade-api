package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	
}

func (h *Handler) wsEndPoint(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(c.Writer,c.Request,nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer ws.Close()

	reader(c, ws)
}

func reader(c *gin.Context, ws *websocket.Conn) {
	for {
		//Read Message from client
		mt, message, err := ws.ReadMessage()
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			break
		}
		//If client message is ping will return pong
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//Response message to client
		err = ws.WriteMessage(mt, message)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			break
		}
	}
}