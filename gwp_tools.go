package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/codegangsta/martini-contrib/render"
	"github.com/fiam/gounidecode/unidecode"
)

func urlize(s string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")

	if err != nil {
		log.Println(err)
	}

	url := reg.ReplaceAllString(unidecode.Unidecode(s), "-")
	url = strings.ToLower(strings.Trim(url, "-"))

	return url
}

func throw404(r render.Render, err error) {
	content := map[string]interface{}{"title": "Not Found", "error": err.Error()}
	r.HTML(404, "404", content)
}
