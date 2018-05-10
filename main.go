// +build !appengine,!appenginevm
package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/wmetaw/echo-on-gae/backend"
)

func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	e.Static("/", "public")

	backend.Routes(e)

	e.Logger.Fatal(e.Start(":8880"))
}
