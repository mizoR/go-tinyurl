package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Tinyurl struct {
	Id   int    `json:"id"`
	Slug string `json:"slug"`
	Url  string `json:"url"`
}

func FindAllTinyurls(db *sql.DB) ([]Tinyurl, error) {
	rows, err := db.Query("SELECT * FROM tinyurls")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tinyurls := []Tinyurl{}

	for rows.Next() {
		tinyurl := Tinyurl{}
		err = rows.Scan(&tinyurl.Id, &tinyurl.Slug, &tinyurl.Url)

		if err != nil {
			return nil, err
		}

		tinyurls = append(tinyurls, tinyurl)
	}

	return tinyurls, nil
}

func FindTinyurlBySlug(db *sql.DB, slug string) (*Tinyurl, error) {
	rows, err := db.Query("SELECT * FROM tinyurls WHERE slug = ?", slug)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	tinyurl := Tinyurl{}

	err = rows.Scan(&tinyurl.Id, &tinyurl.Slug, &tinyurl.Url)

	if err != nil {
		return nil, err
	}

	return &tinyurl, nil
}

func FindTinyurl(db *sql.DB, id int64) (*Tinyurl, error) {
	rows, err := db.Query("SELECT * FROM tinyurls WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tinyurl := Tinyurl{}

	if !rows.Next() {
		return nil, nil
	}

	err = rows.Scan(&tinyurl.Id, &tinyurl.Slug, &tinyurl.Url)

	if err != nil {
		return nil, err
	}

	return &tinyurl, nil
}

func CreateTinyurl(db *sql.DB, slug string, url string) (*Tinyurl, error) {
	result, err := db.Exec(`INSERT INTO "tinyurls" ("slug", "url") VALUES (?, ?)`, slug, url)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	return FindTinyurl(db, id)
}
