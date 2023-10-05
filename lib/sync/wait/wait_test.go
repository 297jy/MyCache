package wait

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	//test1()
	for i := 0; i < 10; i++ {
		testWait()
		fmt.Printf("当前协程的数量：%d\n", runtime.NumGoroutine())
	}
	//time.Sleep(time.Duration(20) * time.Second)

	//time.Sleep(time.Duration(20) * time.Second)

}

func testWait() {
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	w := &Wait{
		wg: sync.WaitGroup{},
	}
	w.Add(1)
	//w.Done()
	//w.Wait()
	//fmt.Printf("当前协程的数量：%d\n", runtime.NumGoroutine())
	w.WaitWithTimeout(time.Duration(2))
}
