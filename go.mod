module github.com/tasavakul/jsonrpcws

go 1.15

require(
    github.com/tasavakul/jsonrpcws/websocket v0.0.0
)

replace(
    github.com/tasavakul/jsonrpcws/websocket => ./websocket
)