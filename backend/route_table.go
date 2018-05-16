package backend

import "github.com/labstack/echo"

// ルーティング
func Routes(e *echo.Echo) {
	e.POST("/items", setItems)
	e.GET("/items", getItems)
	e.GET("/items/tax", getItemTaxes)
	e.GET("/items/:id", getItem)

	e.POST("/users", createUser)
	e.GET("/users", getUsers)
	e.GET("/users/:id", getUser)
	e.GET("/users/logtest", logtest)
	e.GET("/spa", spa)

	e.POST("/task", createTask)
	e.POST("/tasks", createTasks)
	e.GET("/task/:id", getTask)
	e.GET("/testlimit", testlimit)
	e.GET("/testtranslimit", testtranslimit)
}
