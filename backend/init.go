// スタンダード環境用
package backend

import (
	"github.com/labstack/echo"
	"net/http"
)

var e = createMux()

func createMux() *echo.Echo {

	e := echo.New()
	http.Handle("/", e)

	//ルート追加
	Routes(e)
	return e
}
