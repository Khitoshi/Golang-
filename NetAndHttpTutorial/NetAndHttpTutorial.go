package main

import (
	"fmt"
	"net/http"
)

type DB struct {
	//db *sql.DB
	Name string
}

func (db *DB) MyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Println("GET")

	case http.MethodPost:
		fmt.Println("POST")

	default:
		fmt.Println("default")
	}

	fmt.Fprintf(w, "Hello, %s!/n", r.RemoteAddr)
	fmt.Fprintf(w, "param:%s/n", db.Name)
}

func main() {
	fmt.Println("start")
	defer fmt.Println("end")

	db := &DB{Name: "test"}
	http.HandleFunc("/", db.MyHandler)
	//localhost:8080
	http.ListenAndServe(":8080", nil)

}
