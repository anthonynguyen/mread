package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func api_main(c echo.Context) error {
	return c.String(http.StatusOK, "You've reached /api")
}
