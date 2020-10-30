package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/tasavakul/jsonrpcws"
)

func main() {
	e := echo.New()
	ws := jsonrpcws.New()
	ws.Start()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/", jsonrpcws.WSConnect)
	e.Logger.Fatal(e.Start(":1323"))
}
