package cmd

import (
	"github.com/spf13/viper"
	"log"
)

func VerifyParamsOrDie() {
	var url = viper.GetString("url")
	var username = viper.GetString("username")
	var password = viper.GetString("password")
	var searchQuery = viper.GetString("query")
	var spreadsheetId = viper.GetString("spreadsheet")
	var email = viper.GetString("client_email")
	var privkey = viper.GetString("private_key")
	var privkeyId = viper.GetString("private_key_id")

	if url == "" {
		log.Fatalln("Unable to find 'url'")
	}
	if username == "" {
		log.Fatalln("Unable to find 'username'")
	}
	if password == "" {
		log.Fatalln("Unable to find 'password'")
	}
	if searchQuery == "" {
		log.Fatalln("Unable to find 'searchQuery'")
	}
	if spreadsheetId == "" {
		log.Fatalln("Unable to find 'spreadsheetId'")
	}
	if email == "" {
		log.Fatalln("Unable to find 'email'")
	}
	if privkey == "" {
		log.Fatalln("Unable to find 'privkey'")
	}
	if privkeyId == "" {
		log.Fatalln("Unable to find 'privkeyId'")
	}
}
