package jsonrpcws

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

// Reserved Error code
var (
	ParseError     = ErrorCode{-32700, "Parse error"}
	InvalidRequest = ErrorCode{-32600, "Invalid request"}
	MethodNotFound = ErrorCode{-32601, "Method not found"}
	InvalidParam   = ErrorCode{-32602, "Invalid param"}
	InternalError  = ErrorCode{-32603, "Internal error"}
)

// ErrorCode struct
type ErrorCode struct {
	Code    int64
	Message string
}

var (
	rpc      *JSONRPCWS
	upgrader = websocket.Upgrader{}
	handlers = make(map[string]func(rpc *JSONRPCWS, cl *Client, rpcReq *JSONRPCRequest) (*JSONRPCResponse, error))
)

// JSONRPCRequest struct
type JSONRPCRequest struct {
	Jsonrpc *string     `json:"jsonrpc"`
	Method  *string     `json:"method"`
	ID      *string     `json:"id,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Client  *Client
}

// JSONRPCResponse struct
type JSONRPCResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
	ID      *string       `json:"id,omitempty"`
}

// JSONRPCError struct
type JSONRPCError struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JSONRPCWS struct
type JSONRPCWS struct {
	processMessage chan *JSONRPCRequest
	clients        map[string]*Client
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
	go client.StartHandler(rpc)

	return nil
}

// New func
func New() *JSONRPCWS {
	rpc = &JSONRPCWS{
		processMessage: make(chan *JSONRPCRequest),
		clients:        make(map[string]*Client),
	}
	return rpc
}

// RegisterHandler func
func (j *JSONRPCWS) RegisterHandler(method string, handler func(rpc *JSONRPCWS, cl *Client, rpcReq *JSONRPCRequest) (*JSONRPCResponse, error)) {
	handlers[method] = handler
}

// Start func
func (j *JSONRPCWS) Start() {
	go func() {
		for {
			select {
			case message := <-j.processMessage:
				println("Message method:", *message.Method)
				if handler, ok := handlers[*message.Method]; ok {
					resp, err := handler(j, message.Client, message)
					if err != nil {
						message.Client.ResponseError(InternalError, nil, message.ID)
						break
					}
					if resp != nil {
						err = message.Client.Conn.WriteJSON(resp)
						if err != nil {
							println(err.Error())
						}
					}
				} else {
					err := message.Client.ResponseError(MethodNotFound, nil, message.ID)
					if err != nil {
						println(err.Error())
					}
				}
				break
			}
		}
	}()
}

// SendMessage func
func (j *JSONRPCWS) SendMessage(toClientID *string, message *JSONRPCRequest) error {
	if client, ok := j.clients[*toClientID]; ok {
		// TODO: Send message to client
		println("Sending message to ", client)
		err := client.Conn.WriteJSON(message)
		if err != nil {
			return nil
		}
	}
	return nil
}

// AddClient func
func (j *JSONRPCWS) AddClient(clientID string, client *Client) error {
	j.clients[clientID] = client
	return nil
}

func getString(val string) *string {
	return &val
}
