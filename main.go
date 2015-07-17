package main

import (
	"database/sql"
	_ "fmt"
	"github.com/codegangsta/martini"
	"github.com/flosch/pongo2"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"html/template"
	"log"
	"net/http"
	_ "reflect"
	"strconv"
	"time"
)

const (
	TTL_CACHE = 60 // 1 minute for all queries
)

var connection *dbr.Connection
var siteConfig SiteConfig
var cache redis.Conn

func main() {

	// Init db connection
	db, err := sql.Open("mysql", "wp@/wp?charset=utf8")
	if err != nil {
		log.Println(err)
	}

	defer db.Close()
	connection = dbr.NewConnection(db, nil)

	cache, err = redis.Dial("tcp", ":6379")
	if err != nil {
		log.Println(err)
	}

	defer cache.Close()

	InitConfig()

	m := martini.Classic()

	static := martini.Static("static",
		martini.StaticOptions{Fallback: "/404.html", SkipLogging: true})
	m.Use(static)

	m.Get("/", HandleIndex)
	m.Get("/page/:page", HandleIndex)
	m.Get("/:year/:month/:day/:postname", HandlePost)
	m.Get("/:postname", HandlePost)

	m.NotFound(throw404)

	m.Run()
}

func HandlePost(w http.ResponseWriter, r *http.Request, p martini.Params) {
	post, err := GetPost(p["postname"])

	if err != nil {
		throw404(w, r)
		return
	}

	options := GetOptions()

	tpl := pongo2.Must(pongo2.FromFile("templates/post.html"))

	ctxt := pongo2.Context{"title": post.PostTitle, "post": post,
		"options": options}

	err = tpl.ExecuteWriter(ctxt, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleIndex(w http.ResponseWriter, r *http.Request, p martini.Params) {
	startTime := time.Now()
	page, err := strconv.Atoi(p["page"])

	if err != nil {
		page = 0
	}

	posts := GetPosts(page)
	options := GetOptions()

	tpl := pongo2.Must(pongo2.FromFile("templates/posts.html"))

	ctxt := pongo2.Context{"title": "List of posts",
		"posts": posts, "options": options, "page": page}

	err = tpl.ExecuteWriter(ctxt, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println(time.Since(startTime))
}

func throw404(w http.ResponseWriter, r *http.Request) {
	tpl := pongo2.Must(pongo2.FromFile("templates/404.html"))
	ctxt := pongo2.Context{}
	err := tpl.ExecuteWriter(ctxt, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
