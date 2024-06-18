package main

import (
	"log"
	basicUrlShortner "url_shortner/basicUrlShortner"
	sqliteContext "url_shortner/sqliteContext"
)

func main() {
	ctx, err := sqliteContext.New("test.db")
	if err != nil {
		log.Fatalln(err.Error())
	}

	errorFromServer := StartServer(basicUrlShortner.New(ctx))
	log.Println(errorFromServer.Error())
}
