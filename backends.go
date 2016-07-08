package main

import (
	"errors"
	"fmt"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"time"
)

type Backend interface {
	Name() string
	Search(string) ([]SearchResult, error)
}

type SearchResult struct {
	ID              string
	Title           string
	Image           string
	Status          string
	Genres          []string
	LastChapterDate int
	Views           int
}

type MangaEden struct {
	List            []interface{}
	LatestRetrieval int64
}

var MANGA_EDEN = struct {
	MAX_AGE   int64
	LIST_URL  string
	IMAGE_URL string
}{
	3600, // 1 hour, in seconds
	"https://www.mangaeden.com/api/list/0/",
	"https://cdn.mangaeden.com/mangasimg/%s",
}

func (m *MangaEden) Name() string {
	return "MangaEden"
}

func (m *MangaEden) RefreshList() {
	var stringData string

	now := time.Now().Unix()
	if now-m.LatestRetrieval > MANGA_EDEN.MAX_AGE {
		log.Warn("MangaEden list too old, getting another")
		if CONFIG.DEBUG {
			log.Warn("Reading from file")

			data, err := ioutil.ReadFile("./mangaeden.json")
			if err != nil {
				log.Error(err)
				return
			}

			stringData = string(data)
		} else {
			log.Warn("Download a list")
		}

		jq := getJQ(stringData)
		list, err := jq.Array("manga")
		if err != nil {
			return
		}
		m.List = list

		log.Success("Set latest to now")
		m.LatestRetrieval = now
	}
}

func (m *MangaEden) MapStatus(status int) string {
	switch status {
	case 0:
		return "Suspended"
	case 1:
		return "Ongoing"
	case 2:
		return "Completed"
	default:
		return "Unknown"
	}
}

func (m *MangaEden) Search(query string) ([]SearchResult, error) {
	results := make([]SearchResult, 0)
	m.RefreshList()

	if len(m.List) == 0 {
		return nil, errors.New("No list to search")
	}

	for _, jsonManga := range m.List {
		var r SearchResult
		manga := jsonq.NewQuery(jsonManga)

		stringData, err := manga.String("t")
		if err != nil {
			continue
		}

		if !fuzzy(query, stringData) {
			continue
		}

		r.Title = stringData

		stringData, err = manga.String("i")
		if err == nil {
			r.ID = stringData
		}

		stringData, err = manga.String("im")
		if err == nil {
			r.Image = fmt.Sprintf(MANGA_EDEN.IMAGE_URL, stringData)
		}

		intData, err := manga.Int("s")
		if err == nil {
			r.Status = m.MapStatus(intData)
		}

		floatData, err := manga.Float("ld")
		if err == nil {
			r.LastChapterDate = int(floatData)
		}

		intData, err = manga.Int("h")
		if err == nil {
			r.Views = intData
		}

		results = append(results, r)
	}

	return results, nil
}
