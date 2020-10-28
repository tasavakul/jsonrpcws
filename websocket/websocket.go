package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	upgrader = websocket.Upgrader{}
	handlers = make(map[string]func(cl *Client))
)

// JSONRPCRequest struct
type JSONRPCRequest struct {
	Jsonrpc *string     `json:"jsonrpc"`
	Method  *string     `json:"method"`
	ID      *string     `json:"id,omitempty"`
	Params  interface{} `json:"params,omitempty"`
}

// Init func
func init() {
	handlers["register"] = register
}

// WSConnect func
func WSConnect(c echo.Context) error {

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := &Client{
		Conn: conn,
	}

	go client.StartHandler()

	return nil
}
