package main

import (
	"log"
	sqlUS "url_shortner/sql_url_shortner"
	// ust "url_shortner/test_url_shortner"
)

func main() {
	errorFromServer := StartServer(sqlUS.New())
	log.Println(errorFromServer.Error())
}
