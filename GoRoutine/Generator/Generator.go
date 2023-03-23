package main

import (
	"fmt"
	"time"
)

// 複数のgorutineを合体させる場合に使用
func fanIn(ch1, ch2 <-chan string) <-chan string {
	new_ch := make(chan string)

	/*//必ずHello\n Bye の順でoutputされる
	go func() {
		for {
			new_ch <- <-ch1
		}
	}()

	go func() {
		for {
			new_ch <- <-ch2
		}
	}()
	*/

	go func() {
		for { //ランダムに選ばれる
			select {
			case s := <-ch1:
				new_ch <- s
			case s := <-ch2:
				new_ch <- s
			}
		}
	}()

	return new_ch
}

// 戻り値の値を　<-chan　にすると受信専用であることを明示できる
func generator(msg string, quit chan string) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; ; i++ {
			select {
			case ch <- fmt.Sprintf("%s %d", msg, i):
				time.Sleep(time.Second)
			case <-quit:
				quit <- "see you"
				return

			}
		}
		//close(ch)
	}()
	return ch
}

func main() {
	fmt.Println("Start\n")
	defer fmt.Println("\nFinish")

	//ch := fanIn(generator("Hello"), generator("Bye"))
	quit := make(chan string)
	ch := generator("Hi!", quit)
	timeout := time.After(10 * time.Second) //合計の処理が10s以上の場合
	for i := 0; i < 10; i++ {
		select {
		case s := <-ch:
			fmt.Println(s)
		case <-timeout:
			//case <-time.After(1 * time.Second)://1回の処理に1s以上かかる場合
			fmt.Println("waited too long")
			return
		}
	}

	//goroutineが終了していない場合quitに値を渡すことで強制的に終了させる
	//終了させれた場合see you が入る
	quit <- "Bye" //byeと入れることでselectのquitに入り"see you"を受信する
	fmt.Printf("Generator says %s\n", <-quit)

}
