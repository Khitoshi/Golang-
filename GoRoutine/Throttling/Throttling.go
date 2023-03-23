package main

import (
	"fmt"
	"sync"
	"time"
)

func downloadJSON(u string) {
	fmt.Println(u)
	time.Sleep(6 * time.Second)
}

func main() {
	fmt.Println("start\n")
	defer fmt.Println("\nend")

	before := time.Now()

	limit := make(chan struct{}, 20)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)

		i := i
		go func() {
			//limitが20より多くなるとになるとブロックする
			limit <- struct{}{}
			defer wg.Done()
			u := fmt.Sprintf("http://example.com/api/users?id=%d", i)
			downloadJSON(u)
			<-limit //limitから抜き出す
		}()
	}
	wg.Wait()
	//計算上 1回の処理:6s スレッド数:20 回数:100 30s
	//一回の処理 * (スレッド数 / 回数)
	//6 * (100 / 20) = 30s
	fmt.Printf("%v\n", time.Since(before))
}
