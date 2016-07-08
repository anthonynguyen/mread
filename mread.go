package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

var BACKENDS []Backend = []Backend{
	new(MangaEden),
}

func main() {
	loadConfig()

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	setupRoutes(e)

	log.Success("Listening on", CONFIG.PORT)
	e.Run(standard.New(":" + CONFIG.PORT))
}
