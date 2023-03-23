package calc

import (
	"time"
)

func Add(a, b int) int {
	//重い処理
	result := a + b

	time.Sleep(3 * time.Second)

	return result
}
