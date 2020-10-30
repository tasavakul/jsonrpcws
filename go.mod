module github.com/tasavakul/jsonrpcws

go 1.15

require (
	github.com/jsonrpcws/websocket v0.0.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
)

replace github.com/jsonrpcws/websocket => ./websocket
