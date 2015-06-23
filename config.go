package main

import (
	"log"
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

	log.Println(options)

	siteConfig.URLFormat = "YMD"
}

func (c SiteConfig) Get(key string) string {
	return "foo"
}
