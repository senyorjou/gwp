package main

import (
	"log"
	"strconv"
	"strings"
)

const (
	POST_FIELDS = "id, post_date, post_author, post_name, post_title, post_content, post_excerpt"
)

func GetPosts() []*Post {
	var posts []*Post
	dbrSess := connection.NewSession(nil)

	_, err := dbrSess.Select(POST_FIELDS).
		From("wp_posts").
		Where("post_type=?", "post").
		Where("post_status=?", "publish").
		LoadStructs(&posts)

	if err != nil {
		log.Println(err.Error())
	} else {
		loadCategories(posts)
	}

	return posts
}

func GetPost(postName string) (Post, error) {
	var post Post
	var posts []*Post
	dbrSess := connection.NewSession(nil)

	err := dbrSess.Select(POST_FIELDS).
		From("wp_posts").
		Where("post_name=?", postName).
		LoadStruct(&post)

	if err != nil {
		log.Println(err.Error())
	} else {
		posts = append(posts, &post)
		loadCategories(posts)
	}

	return post, err
}

// loadCategories retrieves all categories related to post list passed at once
// it builds an intermediate map to map categories to a single post id
func loadCategories(posts []*Post) {
	type TaxonomiesPosts struct {
		ObjectId  	int64
		Name 		string
		Slug    	string
	}

	var taxposts []*TaxonomiesPosts
	var ids []string
	taxs := map[int64][]Taxs{}
	dbrSess := connection.NewSession(nil)

	for _, post := range(posts) {
		ids = append(ids, strconv.FormatInt(int64(post.Id), 10))
	}
	idstr := strings.Join(ids, ",")

	raw_sql := `SELECT WTR.object_id, WT.name, WT.slug FROM wp_terms WT
		LEFT JOIN wp_term_taxonomy WTT on (WT.term_id = WTT.term_taxonomy_id)
		LEFT JOIN wp_term_relationships WTR on (WT.term_id = WTR.term_taxonomy_id)
		WHERE  WTT.taxonomy = "category" AND WTR.object_id IN (` + idstr + `)`

	_, err := dbrSess.SelectBySql(raw_sql).LoadStructs(&taxposts)
	if err != nil {
		log.Println(err.Error())
	}

	for _, tp:= range(taxposts) {
		taxs[tp.ObjectId] = append(taxs[tp.ObjectId], Taxs{tp.Name, tp.Slug})
	}

	for _, post:= range(posts) {
		post.Categories = append(post.Categories, taxs[post.Id]...)
	}
}

