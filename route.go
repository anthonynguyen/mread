package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func setupRoutes(e *echo.Echo) {
	e.GET("/", route_main)
	e.GET("/search/:query", route_search)
}

func route_main(c echo.Context) error {
	return c.String(http.StatusOK, "You've reached /")
}

func route_search(c echo.Context) error {
	allResults := make(map[string][]SearchResult)
	query := c.Param("query")

	if len(query) < 5 {
		return c.String(http.StatusBadRequest, "Search query is too short")
	}

	for _, b := range BACKENDS {
		results, err := b.Search(query)
		if err != nil {
			log.Error(err)
			continue
		}
		allResults[b.Name()] = results
	}

	return c.JSON(http.StatusOK, allResults)
}
