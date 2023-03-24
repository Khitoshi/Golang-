package fuzzing_test

import (
	"testing"
)

/*
fuzzingテストとは開発者が準備したテストの入力データではなく，
開発者が予期しないであろうランダムで無効なデータを用いてテストする手法
*/
func FuzzDoSomething(f *testing.F) {
	f.Add(3, "test&&&")

	f.Fuzz(func(f *testing.T, i int, s string) {
		doSomething(i, s)
	})

}
