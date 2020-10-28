package websocket

import (
	ws "github.com/gorilla/websocket"
	"github.com/stellar/go/support/log"
)

// Client struct
type Client struct {
	ID   *string `json:"id,omitempty"`
	Conn *ws.Conn
}

// StartHandler func
func (cl *Client) StartHandler() {
	defer func() {
		cl.Conn.Close()
	}()

	for {
		var rpcReq JSONRPCRequest
		err := cl.Conn.ReadJSON(&rpcReq)
		if err != nil {
			log.Errorf(err.Error())
			return
		}

		println("Message Received:", *rpcReq.Method)
	}
}
