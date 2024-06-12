package main

import (
	"fmt"
	"log"
)

func main() {
	errorFromServer := StartServer(NewTestUrlShortner())
	log.Println(errorFromServer.Error())
}

func NewTestUrlShortner() *TestUrlShortner {
	return &TestUrlShortner{
		database: make(map[string]string),
	}
}

type UrlNotFoundError struct{}

func (err UrlNotFoundError) Error() string {
	return "Url not found in system."
}

type TestUrlShortner struct {
	database  map[string]string
	lastIndex int
}

func (urlShortner *TestUrlShortner) SaveUrl(url string) (encriptedUrl string) {
	code := fmt.Sprintf("%d", urlShortner.lastIndex)
	urlShortner.lastIndex++
	urlShortner.database[code] = url
	return code
}

func (urlShortner *TestUrlShortner) GetUrl(code string) (string, error) {
	url, ok := urlShortner.database[code]
	if ok {
		return url, nil
	} else {
		return "", UrlNotFoundError{}
	}
}
