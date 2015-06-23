package main

import (
	"log"
)

func GetPosts() []*Post {
	dbrSess := connection.NewSession(nil)

	var posts []*Post

	_, err := dbrSess.Select("id, post_date, post_author, post_name, post_title, post_content").
		From("wp_posts").
		Where("post_type=?", "post").
		Where("post_status=?", "publish").
		LoadStructs(&posts)

	if err != nil {
		log.Println(err.Error())
	}

	return posts
}
