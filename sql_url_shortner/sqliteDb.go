package SqlUrlShortner

import (
	"database/sql"
	"errors"
	"time"
	"url_shortner/db"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteContext struct {
	connectionString string
}

func NewSqliteCtx(connectionString string) (dbContext *SqliteContext, err error) {
	shortner := &SqliteContext{
		connectionString: connectionString,
	}
	return shortner, shortner.initDatabase()
}

func (sqlUrlShortner *SqliteContext) initDatabase() (err error) {
	db, err := sqlUrlShortner.createDbConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id				INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			url				TEXT NOT NULL,
			lastTime	DATE
		);
	`)

	return err
}

func (sqliteContext *SqliteContext) createDbConnection() (db *sql.DB, err error) {
	return sql.Open("sqlite3", sqliteContext.connectionString)
}

func (sqliteContext *SqliteContext) SaveUrl(url string) (id int64, err error) {
	dbConn, err := sqliteContext.createDbConnection()
	if err != nil {
		return -1, err
	}
	defer dbConn.Close()

	result := dbConn.QueryRow(`
		SELECT *
		FROM urls
		WHERE url == ?;
	`, url)

	var existingRow db.URLRow
	errAtScan := result.Scan(
		&existingRow.Id,
		&existingRow.Url,
		&existingRow.LastTime,
	)

	if errAtScan != nil {
		result, err := dbConn.Exec(`INSERT INTO urls(url, lastTime) VALUES (?, ?);`, url, time.Now())
		if err != nil {
			return -1, err
		}
		return result.LastInsertId()
	}

	return existingRow.Id, nil
}

func (sqliteContext *SqliteContext) GetRow(id int64) (row *db.URLRow, err error) {
	dbConn, err := sqliteContext.createDbConnection()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	result := dbConn.QueryRow(`
		SELECT id, url, lastTime
		FROM urls
		WHERE id = ?;
	`, id)

	var urlRow db.URLRow
	errAtScan := result.Scan(
		&urlRow.Id,
		&urlRow.Url,
		&urlRow.LastTime,
	)
	if errAtScan != nil {
		return nil, errors.New("url not found")
	}

	return &urlRow, nil
}
