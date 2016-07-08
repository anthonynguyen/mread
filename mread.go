package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"log"
)

var BACKENDS []Backend = []Backend{
	MangaEden{},
}

func setupRoutes(e *echo.Echo) {
	e.GET("/", root_main)

	api := e.Group("/api")
	api.GET("", api_main)
}

func main() {
	loadConfig()

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	setupRoutes(e)

	log.Println("Listening on", CONFIG.PORT)
	e.Run(standard.New(":" + CONFIG.PORT))
}
