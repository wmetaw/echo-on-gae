// +build appenginevm

package backend

import (
	"net/http"

	"fmt"
	"github.com/labstack/echo"
	"google.golang.org/appengine"
)

//func createMux() *echo.Echo {
//	e := echo.New()
//	// note: we don't need to provide the middleware or static handlers
//	// for the appengine vm version - that's taken care of by the platform
//	return e
//}

func createMux() *echo.Echo {
	e := echo.New()
	// note: we don't need to provide the middleware or static handlers, that's taken care of by the platform
	// app engine has it's own "main" wrapper - we just need to hook echo into the default handler
	http.Handle("/", e)
	return e
}

func main() {
	// the appengine package provides a convenient method to handle the health-check requests
	// and also run the app on the correct port. We just need to add Echo to the default handler
	//e := echo.New(":8080")
	//http.Handle("/", e)

	appengine.Main()
}
