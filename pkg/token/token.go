package token

import (
	"flag"
	"log"
)

func MustToken() string {
	token := flag.String("bot-token", "", "token for access to telegram bot")
	flag.Parse()
	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}
