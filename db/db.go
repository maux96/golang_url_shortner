package db

import "time"

type DbContext interface {
	SaveUrl(url string) (code string, err error)
	// GetRow(id int64) (row *URLRow, err error)
	GetRowFromCode(code string) (row *URLRow, err error)
}

type URLRow struct {
	Id       int64
	Code     string
	Url      string
	LastTime time.Time
}
