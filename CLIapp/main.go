package main

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/japanese"

	"github.com/PuerkitoBio/goquery"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	AuthorID string
	Author   string
	TitleID  string
	Title    string
	InfoURL  string
	ZipURL   string
}

// 作者とZIPファイルのURLを取得
func findAuthorAndZIP(siteURL string) (string, string) {
	log.Println("query", siteURL)
	//ドキュメント取得
	doc, err := goquery.NewDocument(siteURL)
	if err != nil {
		return "", ""
	}

	//セレクタを取得
	author := doc.Find("table[summary=作家データ]:nth-child(1) tr:nth-child(2) td:nth-child(2)").Text()
	zipURL := ""

	//セレクタの取得してセレクタをある分だけ回す
	doc.Find("table.download a").Each(func(n int, elem *goquery.Selection) {

		//属性を取得
		href := elem.AttrOr("href", "")
		//属性の最後が.zipである場合
		if strings.HasSuffix(href, ".zip") {
			zipURL = href
		}
	})

	if zipURL == "" {
		return author, ""
	}

	//zipURLが http:// 又は https:// から始まっていればtrue
	if strings.HasPrefix(zipURL, "http://") || strings.HasPrefix(zipURL, "https://") {
		return author, zipURL
	}

	//siteURLを要素ごとに分解
	u, err := url.Parse(siteURL)
	if err != nil {
		return author, ""
	}
	//zipURLとpathを結合させる
	u.Path = path.Join(path.Dir(u.Path), zipURL)

	return author, u.String()
}

var pageURLFormat = `https://www.aozora.gr.jp/cards/%s/card%s.html`

// URLからDOMオブジェクトを作成
func findEntries(siteURL string) ([]Entry, error) {
	//ドキュメント取得
	doc, err := goquery.NewDocument(siteURL)
	if err != nil {
		return nil, err
	}

	entries := []Entry{}

	//正規化
	pat := regexp.MustCompile(`.*/cards/([0-9]+)/card([0-9]+).html$`)
	//ドキュメントからセレクタの中のリンクURLを取得してリンクURLをある分だけ回す
	doc.Find("ol li a").Each(func(index int, elem *goquery.Selection) {
		//正規化する
		token := pat.FindStringSubmatch(elem.AttrOr("href", ""))
		if len(token) != 3 {
			return
		}
		title := elem.Text()
		//トークンとURLを組み合わせページURLにする
		pageURL := fmt.Sprintf(pageURLFormat, token[1], token[2])

		//
		author, zipURL := findAuthorAndZIP(pageURL)
		println(zipURL)

		if zipURL != "" {
			entries = append(entries, Entry{
				AuthorID: token[1],
				Author:   author,
				TitleID:  token[2],
				Title:    title,
				InfoURL:  siteURL,
				ZipURL:   zipURL,
			})
		}

	})

	return entries, nil
}

func extractText(zipURL string) (string, error) {
	//URL情報を取得
	resp, err := http.Get(zipURL)
	if err != nil {
		return "", err
	}
	//deferで先に閉じる処理をする
	defer resp.Body.Close()

	//取得下URLの内容を読み込む
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//変数bにダウンロードしたZIPファイルを読み込む
	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return "", err
	}

	//ファイル一覧をloop
	for _, file := range r.File {
		//pathから拡張子を取得し判別
		if path.Ext(file.Name) == ".txt" {
			//ファイルを開く
			f, err := file.Open()
			if err != nil {
				return "", err
			}
			//ファイル読み込み
			b, err := ioutil.ReadAll(f)
			f.Close()
			if err != nil {
				return "", err
			}
			//文字列をUTF-8に変換
			b, err = japanese.ShiftJIS.NewDecoder().Bytes(b)
			if err != nil {
				return "", err
			}

			return string(b), nil
		}
	}

	return "", errors.New("nontents not found")
}

// データベース登録処理
//
//dsn:sqliteファイル名
func setupDB(dsn string) (*sql.DB, error) {
	//dbファイルを開く
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	//dbを閉じる
	defer db.Close()

	//SQLクエリ
	queries := []string{
		`CREATE TABLE IF NOT EXISTS authors(author_id TEXT, author TEXT, PRIMARY KEY (author_id))`,
		`CREATE TABLE IF NOT EXISTS contents(author_id TEXT, title_id TEXT, title TEXT, content TEXT, PRIMARY KEY (author_id, title_id))`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS contents_fts USING fts4(words)`,
	}

	//SQLクエリを実行
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func addEntry(db *sql.DB, entry *Entry, content string) error {
	//SQLクエリ実行
	_, err := db.Exec(
		`REPLACE INTO authors(author_id, author) VALUES(?, ?)`,
		entry.AuthorID,
		entry.Author,
	)
	if err != nil {
		return nil
	}

	//
	res, err := db.Exec(
		`REPLACE INTO contents(author_id, titile_id, title, content) VALUES(?, ?, ?, ?)`,
		entry.AuthorID,
		entry.TitleID,
		entry.Title,
		content,
	)
	if err != nil {
		return nil
	}

	docID, err := res.LastInsertId()
	if err != nil {
		return nil
	}

	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil
	}

	seg := t.Wakati(content)
	_, err = db.Exec(
		`REPLACE INTO contents_fts(docid, words) VALUES(?, ?)`,
		docID,
		strings.Join(seg, " "),
	)
	if err != nil {
		return nil
	}
	return nil
}

func main() {
	db, err := setupDB("database.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	listURL := "https://www.aozora.gr.jp/index_pages/person879.html"

	entries, err := findEntries(listURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("found %d entries\n", len(entries))

	for _, entry := range entries {
		log.Printf("adding %+v\n", entry)
		content, err := extractText(entry.ZipURL)
		if err != nil {
			log.Println(err)
			continue
		}

		err = addEntry(db, &entry, content)
		if err != nil {
			log.Println(err)
			continue
		}

		//fmt.Println(entry.InfoURL)
		//fmt.Println(content)
	}

	fmt.Println("\nend")
}
