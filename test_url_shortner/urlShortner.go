package TestUrlShortner

import "fmt"

type UrlNotFoundError struct{}

func (err UrlNotFoundError) Error() string {
	return "Url not found in system."
}

type TestUrlShortner struct {
	database  map[string]string
	lastIndex int
}

func New() *TestUrlShortner {
	return &TestUrlShortner{
		database: make(map[string]string),
	}
}

func (urlShortner *TestUrlShortner) SaveUrl(url string) (encriptedUrl string, err error) {
	code := fmt.Sprintf("%d", urlShortner.lastIndex)
	urlShortner.lastIndex++
	urlShortner.database[code] = url
	return code, nil
}

func (urlShortner *TestUrlShortner) GetUrl(code string) (string, error) {
	url, ok := urlShortner.database[code]
	if ok {
		return url, nil
	} else {
		return "", UrlNotFoundError{}
	}
}
