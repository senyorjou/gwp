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
	"time"
	_ "net/http"
	_ "reflect"
	"strconv"
)

var connection *dbr.Connection
var siteConfig SiteConfig

func main() {

	// Init db connection
	db, err := sql.Open("mysql", "wp@/wp?charset=utf8")
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

	m.Get("/", HandleIndex)
	m.Get("/page/:page", HandleIndex)
	m.Get("/:year/:month/:day/:postname", HandlePost)
	m.Get("/:postname", HandlePost)

	m.NotFound(throw404)

	m.Run()
}

func HandleIndex(r render.Render, p martini.Params) {
	startTime := time.Now()
	page, err := strconv.Atoi(p["page"])
	if (err != nil) {
		page = 0
	}

	posts := GetPosts(page)
	options := GetOptions()

	content := map[string]interface{}{"title": "List of posts",
		"posts": posts, "options": options, "elapsed": time.Since(startTime),
		"page": page}

	r.HTML(200, "posts", content)
}

func HandlePost(r render.Render, p martini.Params) {
	post, err := GetPost(p["postname"])
	if err != nil {
		throw404(r)
		return
	}
	content := map[string]interface{}{"title": post.PostTitle,
		"post": post}
	r.HTML(200, "post", content)
}

func throw404(r render.Render) {
	content := map[string]interface{}{"title": "Not Found"}
	r.HTML(404, "404", content)
}
