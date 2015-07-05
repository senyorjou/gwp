package main

import (
	_ "log"
	"strings"
)

type Taxs struct {
	Name     string
	Slug     string
	Taxonomy string
	Count    int
}

type Post struct {
	Id          int64
	PostDate    string
	PostAuthor  string
	PostName    string
	PostTitle   string
	PostContent string
	PostExcerpt string
	PostType    string
	MenuOrder   int // only for pages
	Categories  []Taxs
	Tags        []Taxs
}

type Options struct {
	Pages      []*Post
	Categories []Taxs
	Tags       []Taxs
}

func (p Post) Permalink() string {
	var url string

	switch siteConfig.URLFormat {
	case "YMD":
		url = p.PostDate[0:10]
	case "YM":
		url = p.PostDate[0:7]
	}

	return strings.Replace(url, "-", "/", 2) + "/" + p.PostName
}
