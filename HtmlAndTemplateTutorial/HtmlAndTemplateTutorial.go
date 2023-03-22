package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type Person struct {
	Name  string
	Age   int
	Email string
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
	//とりあえず 今回はここで情報を設定(dbを作るのがめんどくさいので)
	person := Person{
		Name:  "张三",
		Age:   18,
		Email: "efpyi@example.com",
	}

	//htmlを取得し解析する
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		panic(err)
	}

	//HTMLテンプレートを実行
	err = tmpl.Execute(w, person)
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Printf("start\n\n")

	http.HandleFunc("/", MyHandler)

	//ウェブサーバを起動
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n\nend")
}
