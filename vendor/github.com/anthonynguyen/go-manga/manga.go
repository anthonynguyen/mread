package manga

import (
	"errors"
	"fmt"
	"sort"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

const MIN_QUERY_LENGTH = 5

type Backend interface {
	Name() string
	Search(string) ([]SearchResult, error)
	Manga(string) (MangaResult, error)
	Chapter(string) ([]string, error)
}

var BACKENDS []Backend = []Backend{
	new(MangaEden),
}

type SearchResult struct {
	ID              string
	Title           string
	Image           string
	Status          string
	Genres          []string
	LastChapterDate string
	Views           int64
}

type ChapterInfo struct {
	Number string
	Date   string
	Title  string
	ID     string
}

type MangaResult struct {
	Title           string
	Image           string
	Status          string
	Genres          []string
	LastChapterDate string
	Views           int64
	Description     string
	NumChapters     int64
	Chapters        []ChapterInfo
}

// So that we can sort
var query string

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

func Search(query string) (map[string][]SearchResult, error) {
	if len(query) < MIN_QUERY_LENGTH {
		return nil, errors.New(fmt.Sprintf("Search query is too short (must be at least %d characters)", MIN_QUERY_LENGTH))
	}

	all := make(map[string][]SearchResult)
	numResults := 0
	for _, b := range BACKENDS {
		results, err := b.Search(query)
		if err != nil {
			continue
		}

		sort.Sort(ByLevenshteinDistance(results))
		all[b.Name()] = results
		numResults += len(results)
	}

	if numResults < 1 {
		return nil, errors.New("No results found")
	}

	return all, nil
}

func Manga(backend string, id string) (MangaResult, error) {
	for _, b := range BACKENDS {
		if backend == b.Name() {
			return b.Manga(id)
		}
	}

	return MangaResult{}, errors.New("Backend not found")
}

func Chapter(backend string, id string) ([]string, error) {
	for _, b := range BACKENDS {
		if backend == b.Name() {
			return b.Chapter(id)
		}
	}

	return nil, errors.New("Backend not found")
}
