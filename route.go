package main

import (
	"github.com/anthonynguyen/go-manga"
	"github.com/labstack/echo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"net/http"
	"sort"
)

type ViewData struct {
	Failed  bool
	Message string
	Query   string
	Backend string
	Data    interface{}
}

// So that we can sort
var query string

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

type ByLevenshteinDistance []manga.SearchResult

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
	allResults := make(map[string][]manga.SearchResult)
	query = c.QueryParam("q")

	if len(query) < 5 {
		return c.Render(http.StatusBadRequest, "search", ViewData{
			Failed:  true,
			Message: "Search query is too short",
			Query:   query,
		})
	}

	numResults := 0
	for _, b := range BACKENDS {
		results, err := b.Search(query)
		if err != nil {
			log.Error(err)
			continue
		}

		sort.Sort(ByLevenshteinDistance(results))
		allResults[b.Name()] = results
		numResults += len(results)
	}

	if numResults < 1 {
		return c.Render(http.StatusOK, "search", ViewData{
			Failed:  true,
			Query:   query,
			Message: "No results found",
		})
	}

	return c.Render(http.StatusOK, "search", ViewData{
		Failed: false,
		Query:  query,
		Data:   allResults,
	})
}

func route_manga(c echo.Context) error {
	requestedBackend := c.Param("backend")
	requestedID := c.Param("id")

	for _, backend := range BACKENDS {
		if requestedBackend == backend.Name() {
			result, err := backend.Manga(requestedID)
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
	}

	return c.Render(http.StatusNotFound, "manga", ViewData{
		Failed:  true,
		Message: "Backend not found",
	})
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
