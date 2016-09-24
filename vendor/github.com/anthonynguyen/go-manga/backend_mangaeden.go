package manga

import (
	"errors"
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/parnurzeal/gorequest"
	"html"
	"io/ioutil"
	"strconv"
	"time"
)

type MangaEden struct {
	List            []*jason.Object
	LatestRetrieval int64
}

var MANGA_EDEN = struct {
	MAX_AGE     int64
	LIST_URL    string
	IMAGE_URL   string
	INFO_URL    string
	CHAPTER_URL string
}{
	3600, // 1 hour, in seconds
	"https://www.mangaeden.com/api/list/0/",
	"https://cdn.mangaeden.com/mangasimg/%s",
	"https://www.mangaeden.com/api/manga/%s/",
	"https://www.mangaeden.com/api/chapter/%s/",
}

func (m *MangaEden) Name() string {
	return "MangaEden"
}

func (m *MangaEden) RefreshList() {
	var data []byte
	var err error

	now := time.Now().Unix()
	if now-m.LatestRetrieval > MANGA_EDEN.MAX_AGE {
		// if CONFIG.DEBUG {
		if false {
			data, err = ioutil.ReadFile("./mangaeden.json")
			if err != nil {
				return
			}
		} else {
			_, body, errs := gorequest.New().Get(MANGA_EDEN.LIST_URL).End()
			if errs != nil {
				return
			}
			data = []byte(body)
		}

		j, err := jason.NewObjectFromBytes(data)
		if err != nil {
			return
		}

		list, err := j.GetObjectArray("manga")
		if err != nil {
			return
		}

		m.List = list

		m.LatestRetrieval = now
	}
}

func (m *MangaEden) MapStatus(status int64) string {
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

	for _, manga := range m.List {
		var r SearchResult

		stringData, err := manga.GetString("t")
		if err != nil {
			continue
		}

		if !fuzzy(query, stringData) {
			continue
		}

		r.Title = stringData

		stringData, err = manga.GetString("i")
		if err == nil {
			r.ID = stringData
		}

		stringData, err = manga.GetString("im")
		if err == nil {
			r.Image = fmt.Sprintf(MANGA_EDEN.IMAGE_URL, stringData)
		}

		intData, err := manga.GetInt64("s")
		if err == nil {
			r.Status = m.MapStatus(intData)
		}

		floatData, err := manga.GetFloat64("ld")
		if err == nil {
			r.LastChapterDate = relTime(floatData)
		}

		intData, err = manga.GetInt64("h")
		if err == nil {
			r.Views = intData
		}

		stringArrayData, err := manga.GetStringArray("c")
		if err == nil {
			r.Genres = stringArrayData
		}

		results = append(results, r)
	}

	return results, nil
}

func (m *MangaEden) Manga(id string) (MangaResult, error) {
	var result MangaResult
	url := fmt.Sprintf(MANGA_EDEN.INFO_URL, id)

	resp, body, errs := gorequest.New().Get(url).End()
	if errs != nil {
		return result, errs[0]
	}

	if resp.StatusCode != 200 {
		return result, errors.New("Status: " + resp.Status)
	}

	data := []byte(body)

	manga, err := jason.NewObjectFromBytes(data)
	if err != nil {
		return result, err
	}

	stringData, err := manga.GetString("title")
	if err == nil {
		result.Title = stringData
	}

	stringData, err = manga.GetString("image")
	if err == nil {
		result.Image = fmt.Sprintf(MANGA_EDEN.IMAGE_URL, stringData)
	}

	intData, err := manga.GetInt64("status")
	if err == nil {
		result.Status = m.MapStatus(intData)
	}

	stringArrayData, err := manga.GetStringArray("categories")
	if err == nil {
		result.Genres = stringArrayData
	}

	floatData, err := manga.GetFloat64("last_chapter_date")
	if err == nil {
		result.LastChapterDate = relTime(floatData)
	}

	intData, err = manga.GetInt64("hits")
	if err == nil {
		result.Views = intData
	}

	stringData, err = manga.GetString("description")
	if err == nil {
		result.Description = html.UnescapeString(stringData)
	}

	intData, err = manga.GetInt64("chapters_len")
	if err == nil {
		result.NumChapters = intData
	}

	valueArrayData, err := manga.GetValueArray("chapters")
	if err == nil {
		for _, chapterJSON := range valueArrayData {
			var chapter ChapterInfo
			arr, err := chapterJSON.Array()
			if err != nil || len(arr) != 4 {
				continue
			}

			floatData, err := arr[0].Float64()
			if err == nil {
				chapter.Number = strconv.FormatFloat(floatData, 'f', -1, 64)
			}

			floatData, err = arr[1].Float64()
			if err == nil {
				chapter.Date = getDate(floatData)
			}

			stringData, err = arr[2].String()
			if err == nil {
				chapter.Title = stringData
			}

			stringData, err = arr[3].String()
			if err != nil {
				// Only the ID is strictly needed
				continue
			}
			chapter.ID = stringData

			result.Chapters = append(result.Chapters, chapter)
		}
	}

	return result, nil
}

// Returns a list of image urls
func (m *MangaEden) Chapter(id string) ([]string, error) {
	result := make([]string, 0)

	url := fmt.Sprintf(MANGA_EDEN.CHAPTER_URL, id)
	resp, body, errs := gorequest.New().Get(url).End()
	if errs != nil {
		return result, errs[0]
	}

	if resp.StatusCode != 200 {
		return result, errors.New("Status: " + resp.Status)
	}

	data := []byte(body)

	chapter, err := jason.NewObjectFromBytes(data)
	if err != nil {
		return result, err
	}

	pages, err := chapter.GetValueArray("images")
	if err == nil {
		// Assume (pray) that the list is in reverse order and just go through it backwards
		for i := len(pages) - 1; i >= 0; i-- {
			arr, err := pages[i].Array()
			if err != nil || len(arr) != 4 {
				// we dead
				continue
			}

			stringData, err := arr[1].String()
			if err == nil {
				result = append(result, fmt.Sprintf(MANGA_EDEN.IMAGE_URL, stringData))
			}
		}

		return result, nil
	}

	return result, err
}
