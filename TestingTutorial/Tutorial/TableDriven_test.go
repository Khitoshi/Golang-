package hsd

import (
	"reflect"
	//"runtime"
	"testing"
)

func TestStringDistance(t *testing.T) {

	tests := []struct {
		name string
		lhs  string
		rhs  string
		want int
	}{
		//テーブルで入力と期待値を纏めておくことでテストの抜けを見つけやすくなる
		//そしてテストの追加も簡単になる
		//テーブルに持たせる項目はおよそ以下になる
		//テストの名前(name)・入力(lhsやrhs)・期待値(want)
		{name: "lhs is longer than rhs", lhs: "foo", rhs: "fo", want: -1},
		{name: "rhs is longer than rhs", lhs: "fo", rhs: "foo", want: -1},
		{name: "No diff", lhs: "foo", rhs: "foo", want: 0},
		{name: "1 diff", lhs: "foo", rhs: "foh", want: 1},
		{name: "2 diff", lhs: "foo", rhs: "fhh", want: 2},
		{name: "3 diff", lhs: "foo", rhs: "hhh", want: 3},
		{name: "multibyte", lhs: "あいえ", rhs: "あいえ", want: 1}, //FAILD
	}

	for _, tc := range tests {
		got := StringDistance(tc.lhs, tc.rhs)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
		}
	}
}

/*
// 特定のOSのスルー
func TestReadData(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.skipf("skipping on %s", runtime.GOOS)
	}

	if runtime.GOOS != "linux" {
		t.skipf("skipping on %s", runtime.GOOS)
	}

	if runtime.GOOS != "mac" {
		t.skipf("skipping on %s", runtime.GOOS)
	}
}
*/
