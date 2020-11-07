package jsonrpcws

import (
	"sync"

	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

// Client struct
type Client struct {
	ID               *string `json:"id,omitempty"`
	Conn             *ws.Conn
	RunningRequestID *int64
	SentRequest      map[string]*JSONRPCRequest
	rpc              *JSONRPCWS
	mu               sync.Mutex
}

// NewClient func
func NewClient(rpc *JSONRPCWS) *Client {
	client := new(Client)
	client.rpc = rpc
	return client
}

// StartHandler func
func (cl *Client) StartHandler(rpc *JSONRPCWS) {
	defer func() {
		cl.Conn.Close()
	}()

	if rpc.OnCloseHandler != nil {
		cl.Conn.SetCloseHandler(func(code int, text string) error {
			if cl.ID != nil {
				return rpc.OnCloseHandler(*cl.ID, code, text)
			}
			return nil
		})
	}

	for {
		var rpcReq *JSONRPCMessage
		err := cl.Conn.ReadJSON(&rpcReq)
		if err != nil {
			println(err.Error())
			return
		}

		// println("Method Received:", *rpcReq.Method)
		// println("Method ID:", *rpcReq.ID)
		PrintJSON(rpcReq)

		rpcReq.Client = cl
		rpc.processMessage <- rpcReq
	}
}

// ResponseError func
func (cl *Client) ResponseError(errorCode JsonrpcError, data interface{}, id *string) error {
	err := cl.rpc.SendResponse(cl, cl.GenerateResponseError(errorCode, data, id))
	if err != nil {
		println(err.Error())
		return err
	}
	return nil
}

// GenerateResponseError func
func (cl *Client) GenerateResponseError(errorCode JsonrpcError, data interface{}, id *string) *JSONRPCResponse {
	return &JSONRPCResponse{
		Jsonrpc: "2.0",
		Error: &JSONRPCError{
			Code:    errorCode.Code,
			Message: errorCode.Message,
			Data:    data,
		},
		ID: id,
	}
}

// GenerateResponseResult func
func (cl *Client) GenerateResponseResult(data interface{}, id *string) *JSONRPCResponse {
	return &JSONRPCResponse{
		Jsonrpc: "2.0",
		Result:  data,
		ID:      id,
	}
}

// NewRequestID func
func (cl *Client) NewRequestID() *string {
	// if cl.RunningRequestID == nil {
	// 	cl.RunningRequestID = GetInt64(1)
	// } else {
	// 	*cl.RunningRequestID++
	// }

	// return cl.RunningRequestID
	return GetString(uuid.New().String())
}

// GetInt64 func
func GetInt64(val int64) *int64 {
	return &val
}

// GetString func
func GetString(val string) *string {
	return &val
}
