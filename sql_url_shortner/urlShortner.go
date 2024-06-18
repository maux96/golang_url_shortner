package SqlUrlShortner

import (
	"log"
	"url_shortner/db"
)

type BasicUrlShortner struct {
	ctx db.DbContext
}

func New() *BasicUrlShortner {
	ctx, err := NewSqliteCtx("test.db")
	if err != nil {
		log.Fatalf("Problem creating the sqlite db context (%s)", err.Error())
	}

	shortner := BasicUrlShortner{
		ctx: ctx,
	}
	return &shortner
}

func (urlShortner *BasicUrlShortner) SaveUrl(url string) (urlCode string, err error) {
	return urlShortner.ctx.SaveUrl(url)
}

func (urlShortner *BasicUrlShortner) GetUrl(code string) (str string, err error) {
	row, err := urlShortner.ctx.GetRowFromCode(code)
	if err != nil {
		return "", err
	}
	return row.Url, nil
}
