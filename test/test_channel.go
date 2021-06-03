package main

import (
	"fmt"
	"time"
)

type Req struct {
	req int
	out chan int
}

type Chan struct {
	req chan Req
}

var c = Chan{
	req: make(chan Req, 100),
}

func calc(x Req) int {
	return x.req * x.req
}
func start() {
	for {
		select {
		case x := <-c.req:
			y := calc(x)
			x.out <- y
			close(x.out)
		}
	}
}
func add(i int) {
	out := make(chan int, 1)
	c.req <- Req{i, out}
	select {
	case x := <-out:
		fmt.Println(x)
	case <-time.After(time.Millisecond * 100):
		fmt.Println("timeout ")
	}
}
func main() {
	N := 100
	go start()
	for i := 0; i < N; i++ {
		add(i)
	}
	time.Sleep(time.Second * 3)
}
