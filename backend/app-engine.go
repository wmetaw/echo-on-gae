// +build appengine

// スタンダード環境用
package backend

import (
	"net/http"

	"github.com/labstack/echo"
)

func createMux() *echo.Echo {
	e := echo.New()
	http.Handle("/", e)
	return e
}
