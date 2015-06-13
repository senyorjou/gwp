package main

import (

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "fmt"

    "github.com/kisielk/sqlstruct"
)

func main() {
    db, err := sql.Open("mysql", "root@/wp?charset=utf8")
    if err !=  nil {
        fmt.Println(err)
    }

    defer db.Close()

    type Post struct {
        ID int
        PostAuthor string `sql:"post_author"`
        PostTitle string `sql:"post_title"`
        PostContent string `sql:"post_content"`
    }

    query := fmt.Sprintf("SELECT %s FROM wp_posts WHERE post_type='%s'",
                         sqlstruct.Columns(Post{}), "post")
    fmt.Println(query)

    rows, err := db.Query(query)
    if err !=  nil {
        fmt.Println(err)
    }

    defer rows.Close()

    var p Post
    for rows.Next() {
        err = sqlstruct.Scan(&p, rows)
        if err !=  nil {
            fmt.Println(err)
        }

        fmt.Println(p.ID, p.PostAuthor, p.PostTitle, p.PostContent)
    }

}