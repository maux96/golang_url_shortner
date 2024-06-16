package SqlUrlShortner

import (
	"fmt"
	"log"
	"strconv"
)

type SqlLiteUrlShortner struct {
	ctx *SqliteContext
}

func New() *SqlLiteUrlShortner {
	ctx, err := NewSqliteCtx("test.db")
	if err != nil {
		log.Fatalf("Problem creating the sqlite db context (%s)", err.Error())
	}

	err = ctx.initDatabase()
	if err != nil {
		log.Fatalf("Problem creating the sqlite db context (%s)", err.Error())
	}

	shortner := SqlLiteUrlShortner{
		ctx: ctx,
	}
	return &shortner
}

func (urlShortner *SqlLiteUrlShortner) SaveUrl(url string) (urlCode string, err error) {
	id, err := urlShortner.ctx.SaveUrl(url)
	return fmt.Sprintf("%d", id), err
}

func (urlShortner *SqlLiteUrlShortner) GetUrl(code string) (str string, err error) {
	id, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		return "", err
	}
	row, err := urlShortner.ctx.GetRow(id)
	if err != nil {
		return "", err
	}
	return row.Url, nil
}
