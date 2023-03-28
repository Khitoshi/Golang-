package main_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	//
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	//dbを閉じる
	defer db.Close()

	query := "虫 AND ココア"
	rows, err := db.Query(`
		SELECT
		a.author,
		c.title,
		FROM
		contents c
		INNER JOIN authors a
		ON a.author_id = c.author_id
		INNER JOIN contents_fts f
		ON c.rowid = f.docid
		AND words MATCH ?
		`, query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var author, title string
		err = rows.Scan(&author, &title)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(author, title)
	}
}
