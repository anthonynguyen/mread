package main

import (
	"errors"
	"github.com/eknkc/amber"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"os"
)

var BACKENDS []Backend = []Backend{
	new(MangaEden),
}

type Views struct {
	templates map[string]*template.Template
}

func (t *Views) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	for templateName, temp := range t.templates {
		if name == templateName {
			return temp.Execute(w, data)
		}
	}

	return errors.New("Template not found")
}

func main() {
	loadConfig()

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	viewInfo := []struct {
		name string
		path string
	}{
		{"index", "./views/index.amber"},
		{"search", "./views/search.amber"},
		{"manga", "./views/manga.amber"},
	}

	t := &Views{make(map[string]*template.Template)}

	compiler := amber.New()
	for _, v := range viewInfo {
		err := compiler.ParseFile(v.path)
		if err != nil {
			log.Error("Error parsing template file: ", err)
			os.Exit(1)
		}
		tpl, err := compiler.Compile()
		if err != nil {
			log.Error("Error compiling template file: ", err)
			os.Exit(1)
		}
		t.templates[v.name] = tpl
	}

	e.SetRenderer(t)
	setupRoutes(e)
	e.Static("/static", "static")

	log.Success("Listening on", CONFIG.PORT)
	e.Run(standard.New(":" + CONFIG.PORT))
}
