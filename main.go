package main

import (
	"database/sql"
	_ "fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"html/template"
	"log"
	_ "net/http"
	_ "reflect"
)

var connection *dbr.Connection
var siteConfig SiteConfig

func main() {
	// Init db connection
	db, err := sql.Open("mysql", "root@/wp?charset=utf8")
	if err != nil {
		log.Println(err)
	}

	defer db.Close()
	connection = dbr.NewConnection(db, nil)

	InitConfig()

	m := martini.Classic()

	static := martini.Static("static",
		martini.StaticOptions{Fallback: "/404.tmpl", SkipLogging: true})
	m.Use(static)

	m.Use(render.Renderer(render.Options{
		Directory: "templates",
		Layout:    "base",
		Funcs: []template.FuncMap{
			{
				"unescaped": func(args ...interface{}) template.HTML {
					return template.HTML(args[0].(string))
				},
			},
		},
	}))

	m.Get("/", func(r render.Render) {
		//fetch all rows
		posts := GetPosts()
		content := map[string]interface{}{"title": "List of posts",
			"posts": posts}

		r.HTML(200, "posts", content)
	})

	m.Get("/:year/:month/:day/:postname", func(r render.Render, p martini.Params) {
		post, err := GetPost(p["postname"])
		if err != nil {
			throw404(r, err)
			return
		}
		content := map[string]interface{}{"title": post.PostTitle,
			"post": post}
		r.HTML(200, "post", content)
	})

	// 404
	m.NotFound(throw404)

	// RUN
	m.Run()
}
