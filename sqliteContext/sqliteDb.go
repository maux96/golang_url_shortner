package SqliteContext

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"url_shortner/db"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteContext struct {
	connectionString string
}

func New(connectionString string) (dbContext *SqliteContext, err error) {
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
			code			CHAR(255) NOT NULL UNIQUE,
			url				TEXT NOT NULL,
			lastTime	DATE
		);
	`)

	return err
}

func (sqliteContext *SqliteContext) createDbConnection() (db *sql.DB, err error) {
	return sql.Open("sqlite3", sqliteContext.connectionString)
}

func (sqliteContext *SqliteContext) SaveUrl(url string) (code string, err error) {
	dbConn, err := sqliteContext.createDbConnection()
	if err != nil {
		return "", err
	}
	defer dbConn.Close()

	result := dbConn.QueryRow(`
		SELECT id, code, url, lastTime  
		FROM urls
		WHERE url == ?;
	`, url)

	var existingRow db.URLRow
	errAtScan := result.Scan(
		&existingRow.Id,
		&existingRow.Code,
		&existingRow.Url,
		&existingRow.LastTime,
	)
	if errAtScan != nil {
		transCtx, err := dbConn.Begin()
		if err != nil {
			return "", errors.New("start transaction failed")
		}

		result, err := transCtx.Exec(`INSERT INTO urls(url, code, lastTime) VALUES (?, ?, ?);`, url, code, time.Now())
		if err != nil {
			rollbackError := transCtx.Rollback()
			return "", errors.Join(err, rollbackError)
		}

		id, err := result.LastInsertId()
		if err != nil {
			rollbackError := transCtx.Rollback()
			return "", errors.Join(err, rollbackError)
		}

		// we set the code as the id for now
		code = fmt.Sprintf("%d", id)
		_, err = transCtx.Exec(`UPDATE urls SET code = ? WHERE id = ?;`, id, code)
		if err != nil {
			rollbackError := transCtx.Rollback()
			return "", errors.Join(err, rollbackError)
		}

		err = transCtx.Commit()
		if err != nil {
			return "", err
		}
		return code, nil
	}

	return existingRow.Code, nil
}

func (sqliteContext *SqliteContext) GetRowFromCode(code string) (row *db.URLRow, err error) {
	dbConn, err := sqliteContext.createDbConnection()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	result := dbConn.QueryRow(`
		SELECT id, url, lastTime
		FROM urls
		WHERE code = ?;
	`, code)

	urlRow := db.URLRow{Code: code}
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
