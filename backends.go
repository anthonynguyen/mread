package main

type Backend interface {
	Name() string
}

type MangaEden struct{}

func (me MangaEden) Name() string {
	return "MangaEden"
}
