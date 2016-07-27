package main

import (
	"github.com/anthonynguyen/go-manga"
	"github.com/labstack/echo"
	"net/http"
)

type ViewData struct {
	Failed  bool
	Message string
	Query   string
	Backend string
	Data    interface{}
}

func setupRoutes(e *echo.Echo) {
	e.GET("/", route_main)
	e.GET("/search", route_search)
	e.GET("/manga/:backend/:id", route_manga)
	e.GET("/chapter/:backend/:id", route_chapter)
}

func route_main(c echo.Context) error {
	return c.Render(http.StatusOK, "index", ViewData{
		Failed:  true,
		Message: "Use the bar above to search",
	})
}

func route_search(c echo.Context) error {
	query := c.QueryParam("q")
	all, err := manga.Search(query)

	if err != nil {
		return c.Render(http.StatusBadRequest, "search", ViewData{
			Failed:  true,
			Message: err.Error(),
			Query:   query,
		})
	}

	return c.Render(http.StatusOK, "search", ViewData{
		Failed: false,
		Query:  query,
		Data:   all,
	})
}

func route_manga(c echo.Context) error {
	requestedBackend := c.Param("backend")
	requestedID := c.Param("id")

	result, err := manga.Manga(requestedBackend, requestedID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "manga", ViewData{
			Failed:  true,
			Message: err.Error(),
			Backend: requestedBackend,
		})
	}

	return c.Render(http.StatusOK, "manga", ViewData{
		Failed:  false,
		Data:    result,
		Backend: requestedBackend,
	})
}

func route_chapter(c echo.Context) error {
	requestedBackend := c.Param("backend")
	requestedID := c.Param("id")

	result, err := manga.Chapter(requestedBackend, requestedID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, result)
}
