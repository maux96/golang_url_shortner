package BasicUrlShortner

import (
	"url_shortner/db"
)

type BasicUrlShortner struct {
	ctx db.DbContext
}

func New(ctx db.DbContext) *BasicUrlShortner {
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
