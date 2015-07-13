package main

import (
	"log"
	"strconv"
)

type SiteConfig struct {
	BlogName        string
	BlogDescription string
	URLFormat       string
	PostxPage       int
}

func InitConfig() {
	type Option struct {
		OptionName  string
		OptionValue string
	}

	var options []*Option

	dbrSess := connection.NewSession(nil)
	_, err := dbrSess.Select("option_name, option_value").
		From("wp_options").
		LoadStructs(&options)

	if err != nil {
		log.Println(err.Error())
	}

	for _, option := range options {
		switch option.OptionName {
		case "blogname":
			siteConfig.BlogName = option.OptionValue

		case "blogdescription":
			siteConfig.BlogDescription = option.OptionValue

		case "posts_per_page":
			siteConfig.PostxPage, _ = strconv.Atoi(option.OptionValue)
		}
	}

	siteConfig.URLFormat = "YMD"
}
