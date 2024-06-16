package db

import "time"

type DbContext interface {
	SaveUrl(url string) (id int, err error)
	GetRow(id int64) (row *URLRow, err error)
}

type URLRow struct {
	Id       int64
	Url      string
	LastTime time.Time
}
