package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// テストデータをtestdataディレクトリに配置してテストを実施する
func TestDoSomething(t *testing.T) {
	fns, err := filepath.Glob("testdata/*.dat")
	if err != nil {
		t.Fatal(err)
	}

	for _, fn := range fns {
		t.Log(fn)
		//入力データを読む
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			t.Fatal(err)
		}
		//関数呼び出し
		got := doSomething(string(b))

		//.datを.outに入れ替えて結果データを読み込む
		b, err = ioutil.ReadFile(fn[:len(fn)-4] + "out")
		if err != nil {
			t.Fatal(err)
		}
		want := string(b)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf(diff)
		}
	}
}
