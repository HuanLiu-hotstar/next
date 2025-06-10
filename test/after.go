package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func get(ctx context.Context, x int) {
	fmt.Println(x)
	select {
	case <-time.After(time.Second * 10):
		fmt.Printf("timeout %d\n", x)
	}
}
func main() {

	ctx := context.Background()
	N := 100
	wg := sync.WaitGroup{}

	for i := 0; i < N; i++ {

		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			get(ctx, x)
		}(i)
	}
	wg.Wait()

}
