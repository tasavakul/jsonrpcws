package jsonrpcws

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

const (
	jsonrpcVersion = "2.0"
)

// Reserved Error code
var (
	ParseError     = JsonrpcError{-32700, "Parse error"}
	InvalidRequest = JsonrpcError{-32600, "Invalid request"}
	MethodNotFound = JsonrpcError{-32601, "Method not found"}
	InvalidParam   = JsonrpcError{-32602, "Invalid param"}
	InternalError  = JsonrpcError{-32603, "Internal error"}
)

// JsonrpcError struct
type JsonrpcError struct {
	Code    int64
	Message string
}

// CommonError struct
type CommonError struct {
	error
	Code    int64
	Message string
}

// Error Declaration for Job
var (
	DuplicateJobError      = CommonError{Code: 100, Message: "Duplicate job entry"}
	ParameterNotFoundError = CommonError{Code: 101, Message: "Parameter not found"}
	InvalidParameterMetric = CommonError{Code: 102, Message: "Parameters of one Job has multiple metric"}
	ClientNotFound         = CommonError{Code: 103, Message: "Client not found"}
)
var (
	rpc              *JSONRPCWS
	upgrader         = websocket.Upgrader{}
	requestHandlers  = make(map[string]func(rpc *JSONRPCWS, cl *Client, rpcMessage *JSONRPCMessage) (*JSONRPCResponse, error))
	responseHandlers = make(map[string]func(rpc *JSONRPCWS, cl *Client, rpcMessage *JSONRPCMessage) error)
)

// JSONRPCRequest struct
type JSONRPCRequest struct {
	Jsonrpc        *string                          `json:"jsonrpc"`
	Method         *string                          `json:"method"`
	ID             *string                          `json:"id,omitempty"`
	Params         interface{}                      `json:"params,omitempty"`
	Client         *Client                          `json:"-"`
	ResponseHandle func(res *JSONRPCResponse) error `json:"-"`
}

// JSONRPCResponse struct
type JSONRPCResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
	ID      *string       `json:"id,omitempty"`
}

// JSONRPCMessage struct
type JSONRPCMessage struct {
	Jsonrpc        *string                          `json:"jsonrpc"`
	Method         *string                          `json:"method,omitempty"`
	ID             *string                          `json:"id,omitempty"`
	Params         interface{}                      `json:"params,omitempty"`
	Result         interface{}                      `json:"result,omitempty"`
	Error          *JSONRPCError                    `json:"error,omitempty"`
	Client         *Client                          `json:"-"`
	ResponseHandle func(res *JSONRPCResponse) error `json:"-"`
}

// JSONRPCError struct
type JSONRPCError struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JSONRPCWS struct
type JSONRPCWS struct {
	processMessage chan *JSONRPCMessage
	clients        map[string]*Client
	OnCloseHandler func(clientID string, code int, text string) error
}

// WSConnect func
func WSConnect(c echo.Context) error {

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	client := &Client{
		Conn: conn,
		rpc:  rpc,
	}
	go client.StartHandler(rpc)

	return nil
}

// New func
func New() *JSONRPCWS {
	rpc = &JSONRPCWS{
		processMessage: make(chan *JSONRPCMessage),
		clients:        make(map[string]*Client),
	}
	return rpc
}

// RegisterRequestHandler func
func (j *JSONRPCWS) RegisterRequestHandler(method string, handler func(rpc *JSONRPCWS, cl *Client, rpcReq *JSONRPCMessage) (*JSONRPCResponse, error)) {
	requestHandlers[method] = handler
}

// Start func
func (j *JSONRPCWS) Start() {
	go func() {
		for {
			select {
			case message := <-j.processMessage:
				if message.Method != nil {
					println("Suppose to be Request message")
					println("Message method:", *message.Method)
					if handler, ok := requestHandlers[*message.Method]; ok {
						resp, err := handler(j, message.Client, message)
						if err != nil {
							message.Client.ResponseError(InternalError, nil, message.ID)
							break
						}
						if resp != nil {
							PrintJSON(resp)
							err = j.SendResponse(message.Client, resp)
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
				} else {
					println("Suppose to be Response Message")
					if req, ok := message.Client.SentRequest[*message.ID]; ok {
						println("Response from Request :")
						PrintJSON(req)
						// if handler, ok := responseHandlers[*req.Method]; ok {
						// 	err := handler(j, message.Client, message)
						var res JSONRPCResponse
						err := Convert(message, &res)
						if err != nil {
							message.Client.ResponseError(InternalError, nil, message.ID)
							break
						}
						println("Converted Res :")
						PrintJSON(res)
						println("Handler:", req.ResponseHandle)
						err = req.ResponseHandle(&res)
						if err != nil {
							message.Client.ResponseError(InternalError, nil, message.ID)
							break
						}

					} else {
						err := message.Client.ResponseError(InvalidRequest, nil, message.ID)
						if err != nil {
							println(err.Error())
						}
					}
				}
				break
			}
		}
	}()
}

// GetClientByID func
func (j *JSONRPCWS) GetClientByID(clientID *string) *Client {

	if client, ok := j.clients[*clientID]; ok {
		return client
	}
	return nil
}

// SendRequest func
func (j *JSONRPCWS) SendRequest(client *Client, request *JSONRPCRequest) error {
	var mess JSONRPCMessage
	err := Convert(request, &mess)
	if err != nil {
		return err
	}
	mess.ResponseHandle = request.ResponseHandle
	return j.SendMessage(client, &mess)
}

// SendResponse func
func (j *JSONRPCWS) SendResponse(client *Client, response *JSONRPCResponse) error {
	var mess JSONRPCMessage
	err := Convert(response, &mess)
	if err != nil {
		return err
	}
	return j.SendMessage(client, &mess)
}

// SendMessage func
func (j *JSONRPCWS) SendMessage(client *Client, message *JSONRPCMessage) error {
	isRequest := false
	if message.Method != nil {
		isRequest = true
	}

	// TODO: Send message to client
	client.mu.Lock()
	if isRequest {
		message.ID = client.NewRequestID()
	}
	message.Jsonrpc = getString(jsonrpcVersion)
	println("Sending message to ", client)
	PrintJSON(message)
	err := client.Conn.WriteJSON(message)
	if err != nil {
		return nil
	}
	if isRequest {
		var req JSONRPCRequest
		err := Convert(message, &req)
		if err != nil {
			return err
		}
		req.ResponseHandle = message.ResponseHandle
		client.SentRequest[*message.ID] = &req
	}

	client.mu.Unlock()

	return nil
}

// AddClient func
func (j *JSONRPCWS) AddClient(clientID string, client *Client) error {
	j.clients[clientID] = client
	client.ID = &clientID
	client.SentRequest = make(map[string]*JSONRPCRequest)
	return nil
}

func getString(val string) *string {
	return &val
}
