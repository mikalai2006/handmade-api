package handler

import (
	"encoding/json"
	"fmt"
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

	reader(h, c, ws)
}

func reader(h *Handler, c *gin.Context, ws *websocket.Conn) {
	for {
		//Read Message from client
		mt, message, err := ws.ReadMessage()
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			break
		}


		var m WsMessage
		if err := json.Unmarshal([]byte(message), &m); err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		//If client message is ping will return pong
		if m.Method == "find" {
			response, err := h.services.Shop.GetAllShops()
			if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			// response := []string{
			// 	"pong",
			// 	"pong",
			// }
			updateJson, err := json.Marshal(response)
			if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			message = updateJson

		}

		//Response message to client
		err = ws.WriteMessage(mt, message)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			break
		}
	}
}

type WsMessage struct {
	Service  string
	Method string
	Data any
}

func (n *WsMessage) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&n.Service, &n.Method, &n.Data}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Notification: %d != %d", g, e)
	}
	return nil
}