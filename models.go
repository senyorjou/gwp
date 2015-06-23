package main

import (
	_ "log"
	"strings"
)

type Post struct {
	Id          int64
	PostDate    string
	PostAuthor  string
	PostName    string
	PostTitle   string
	PostContent string
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
