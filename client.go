package jsonrpcws

import (
	ws "github.com/gorilla/websocket"
)

// Client struct
type Client struct {
	ID   *string `json:"id,omitempty"`
	Conn *ws.Conn
}

// StartHandler func
func (cl *Client) StartHandler(rpc *JSONRPCWS) {
	defer func() {
		cl.Conn.Close()
	}()

	for {
		var rpcReq *JSONRPCRequest
		err := cl.Conn.ReadJSON(&rpcReq)
		if err != nil {
			println(err.Error())
			return
		}

		println("Method Received:", *rpcReq.Method)
		println("Method ID:", *rpcReq.ID)

		rpcReq.Client = cl
		rpc.processMessage <- rpcReq
	}
}

// ResponseError func
func (cl *Client) ResponseError(errorCode ErrorCode, data interface{}, id *string) error {
	err := cl.Conn.WriteJSON(cl.GenerateResponseError(errorCode, data, id))
	if err != nil {
		println(err.Error())
		return err
	}
	return nil
}

// GenerateResponseError func
func (cl *Client) GenerateResponseError(errorCode ErrorCode, data interface{}, id *string) *JSONRPCResponse {
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

// GenerateReponseResult func
func (cl *Client) GenerateReponseResult(data interface{}, id *string) *JSONRPCResponse {
	return &JSONRPCResponse{
		Jsonrpc: "2.0",
		Result:  data,
		ID:      id,
	}
}
