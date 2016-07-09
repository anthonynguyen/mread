package main

import (
	"github.com/labstack/echo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"net/http"
	"sort"
)

// So that we can sort
var query string

func setupRoutes(e *echo.Echo) {
	e.GET("/", route_main)
	e.GET("/search/:query", route_search)
	e.GET("/manga/:backend/:id", route_manga)
	e.GET("/chapter/:backend/:id", route_chapter)
}

func route_main(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

type ByLevenshteinDistance []SearchResult

func (r ByLevenshteinDistance) Len() int {
	return len(r)
}

func (r ByLevenshteinDistance) Swap(i int, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ByLevenshteinDistance) Less(i int, j int) bool {
	return levenshtein.DistanceForStrings([]rune(query), []rune(r[i].Title), levenshtein.DefaultOptions) < levenshtein.DistanceForStrings([]rune(query), []rune(r[j].Title), levenshtein.DefaultOptions)
}

func route_search(c echo.Context) error {
	allResults := make(map[string][]SearchResult)
	query = c.Param("query")

	if len(query) < 5 {
		return c.String(http.StatusBadRequest, "Search query is too short")
	}

	for _, b := range BACKENDS {
		results, err := b.Search(query)
		if err != nil {
			log.Error(err)
			continue
		}

		sort.Sort(ByLevenshteinDistance(results))

		allResults[b.Name()] = results
	}

	return c.JSON(http.StatusOK, allResults)
}

func route_manga(c echo.Context) error {
	requestedBackend := c.Param("backend")
	requestedID := c.Param("id")

	for _, backend := range BACKENDS {
		if requestedBackend == backend.Name() {
			result, err := backend.Manga(requestedID)
			if err != nil {
				return c.String(http.StatusInternalServerError, "")
			}

			return c.JSON(http.StatusOK, result)
		}
	}

	return c.String(http.StatusNotFound, "Backend not found")
}

func route_chapter(c echo.Context) error {
	requestedBackend := c.Param("backend")
	requestedID := c.Param("id")

	for _, backend := range BACKENDS {
		if requestedBackend == backend.Name() {
			result, err := backend.Chapter(requestedID)
			if err != nil {
				return c.String(http.StatusInternalServerError, "")
			}

			return c.JSON(http.StatusOK, result)
		}
	}

	return c.String(http.StatusNotFound, "Backend not found")
}
