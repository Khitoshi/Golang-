package something

import (
	"fmt"
	"time"
)

func DoSomething() {
	//処理
	time.Sleep(5 * time.Second)
}

func makeSomething(n int) []string {
	r := make([]string, n, n) //new
	//var r []string//old
	for i := 0; i < n; i++ {
		r = append(r, fmt.Sprintf("%05d 何か", i))
	}
	return r
}
