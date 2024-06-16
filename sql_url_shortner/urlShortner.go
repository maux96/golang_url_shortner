package SqlUrlShortner

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqlUrlShortner struct {
	connectionString string
}

type URLRow struct {
	id       string
	url      string
	lastTime time.Time
}

func New(connectionString string) *SqlUrlShortner {
	shortner := &SqlUrlShortner{
		connectionString: connectionString,
	}
	err := shortner.initDatabase()
	if err != nil {
		log.Fatalf("ERROR inicializando la bd. (%s)", err.Error())
	}

	return shortner
}

func (sqlUrlShortner *SqlUrlShortner) createDbConnection() (db *sql.DB, err error) {
	return sql.Open("sqlite3", sqlUrlShortner.connectionString)
}

func (sqlUrlShortner *SqlUrlShortner) initDatabase() (err error) {
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

func (urlShortner *SqlUrlShortner) SaveUrl(url string) (encriptedUrl string, err error) {
	dbConn, err := urlShortner.createDbConnection()
	if err != nil {
		return "", err
	}
	defer dbConn.Close()

	result := dbConn.QueryRow(`
		SELECT *
		FROM urls
		WHERE url == ?;
	`, url)

	var existingRow URLRow
	errAtScan := result.Scan(
		&existingRow.id,
		&existingRow.url,
		&existingRow.lastTime,
	)

	if errAtScan != nil {
		result, err := dbConn.Exec(`INSERT INTO urls(url, lastTime) VALUES (?, ?);`, url, time.Now())
		if err != nil {
			return "", err
		}
		id, err := result.LastInsertId()
		return fmt.Sprintf("%d", id), err
	}

	return existingRow.id, nil
}

func (urlShortner *SqlUrlShortner) GetUrl(code string) (string, error) {
	dbConn, err := urlShortner.createDbConnection()
	if err != nil {
		return "", err
	}
	defer dbConn.Close()

	result := dbConn.QueryRow(`
		SELECT url
		FROM urls
		WHERE id == ?;
	`, code)

	var url string
	errAtScan := result.Scan(
		&url,
	)
	if errAtScan != nil {
		return "", errors.New("url not found")
	}

	return url, nil
}
