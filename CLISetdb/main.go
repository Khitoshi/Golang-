package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/encoding/japanese"
	//"github.com/ikawaha/kagome-dict/ipa"
	//"github.com/ikawaha/kagome/v2/tokenizer"
)

func main() {
	//
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	//dbを閉じる
	defer db.Close()

	//sqlコマンド　記述
	queries := []string{
		` CREATE TABLE IF NOT EXISTS authors( author_id TEXT, author TEXT, PRIMARY KEY (author_id))`,
		` CREATE TABLE IF NOT EXISTS contents( author_id TEXT, title_id TEXT, title TEXT, content TEXT, PRIMARY KEY (author_id, title_id))`,
		` CREATE VIRTUAL TABLE IF NOT EXISTS contents_fts USING fts4(words)`,
	}
	//クエリを発行
	for _, querie := range queries {
		_, err := db.Exec(querie)
		if err != nil {
			log.Fatal(err)
		}
	}
	b, err := os.ReadFile("ababababa.txt")
	if err != nil {
		log.Fatal(err)
	}

	b, err = japanese.ShiftJIS.NewDecoder().Bytes(b)
	if err != nil {
		log.Fatal(err)
	}
	content := string(b)

	res, err := db.Exec(
		`INSERT INTO contents(author_id, title_id, title,content) values(?, ?, ?, ?)`,
		"000879",
		"14",
		"あばばばば",
		content,
	)
	if err != nil {
		log.Fatal(err)
	}

	docID, err := res.LastInsertId()

	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		log.Fatal(err)
	}

	seg := t.Wakati(content)
	_, err = db.Exec(`
	    INSERT INTO contents_fts(docid,words) VALUES(?,?)
	`,
		docID,
		strings.Join(seg, " "),
	)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("end")
}
