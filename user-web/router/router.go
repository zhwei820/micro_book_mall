package router

import (
	"github.com/labstack/echo"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/user-web/handler"
)

func Router() *echo.Echo {
	// Echo instance
	e := echo.New()
	//e.Use(middleware.RequestID())
	e.GET("/", handler.Login)
	e.GET("/home", handler.Home)
	e.GET("/info", handler.Info)
	return e
}
