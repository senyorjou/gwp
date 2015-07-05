package main

import (
	"log"
	"strconv"
	"strings"
)

const (
	POST_FIELDS = "id, post_date, post_author, post_name, post_title, post_content, post_excerpt, menu_order"
)

func GetOptions() Options {
	dbrSess := connection.NewSession(nil)
	var options Options

	_, err := dbrSess.Select(POST_FIELDS).
		From("wp_posts").
		Where("post_type=?", "page").
		Where("post_status=?", "publish").
		OrderBy("menu_order ASC").
		LoadStructs(&options.Pages)

	if err != nil {
		log.Println(err.Error())
	}

	var taxs []*Taxs

	raw_sql := `SELECT WT.name, WT.slug, WTT.taxonomy, WTT.count FROM wp_terms WT
		LEFT JOIN wp_term_taxonomy WTT ON (WTT.term_taxonomy_id = WT.term_id)
		WHERE WTT.taxonomy IN ("category", "post_tag")`

	_, err = dbrSess.SelectBySql(raw_sql).LoadStructs(&taxs)
	if err != nil {
		log.Println(err.Error())
	}

	for _, tax := range taxs {
		switch tax.Taxonomy {
		case "category":
			options.Categories = append(options.Categories,
				Taxs{tax.Name, tax.Slug, tax.Taxonomy, tax.Count})

		case "post_tag":
			options.Tags = append(options.Tags,
				Taxs{tax.Name, tax.Slug, tax.Taxonomy, tax.Count})
		}
	}

	return options
}

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
		ObjectId int64
		Name     string
		Slug     string
		Taxonomy string
		Count    int
	}

	var taxposts []*TaxonomiesPosts
	var ids []string
	cats := map[int64][]Taxs{}
	tags := map[int64][]Taxs{}
	var postType = make(map[int64]string)
	// taxtype := map[int64]string
	dbrSess := connection.NewSession(nil)

	for _, post := range posts {
		ids = append(ids, strconv.FormatInt(int64(post.Id), 10))
	}
	idstr := strings.Join(ids, ",")

	raw_sql := `SELECT WTR.object_id, WT.name, WT.slug, WTT.taxonomy, WTT.count
		FROM wp_terms WT
		LEFT JOIN wp_term_taxonomy WTT on (WT.term_id = WTT.term_taxonomy_id)
		LEFT JOIN wp_term_relationships WTR on (WT.term_id = WTR.term_taxonomy_id)
		WHERE WTR.object_id IN (` + idstr + `)`

	_, err := dbrSess.SelectBySql(raw_sql).LoadStructs(&taxposts)
	if err != nil {
		log.Println(err.Error())
	}

	for _, tp := range taxposts {
		switch tp.Taxonomy {
		case "category":
			cats[tp.ObjectId] = append(cats[tp.ObjectId],
				Taxs{tp.Name, tp.Slug, tp.Taxonomy, tp.Count})

		case "post_tag":
			tags[tp.ObjectId] = append(tags[tp.ObjectId],
				Taxs{tp.Name, tp.Slug, tp.Taxonomy, tp.Count})

		case "post_format":
			postType[tp.ObjectId] = tp.Slug
		}
	}

	for _, post := range posts {
		post.Categories = append(post.Categories, cats[post.Id]...)
		post.Tags = append(post.Tags, tags[post.Id]...)
		post.PostType = postType[post.Id]
	}
}
