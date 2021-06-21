package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func add(ch chan string) string {
	x := ""
	// d := time.After()
	d := time.NewTimer(time.Second * 60)
	defer d.Stop()
	select {
	case x := <-ch:
		fmt.Println(x)
		return x
	case <-d.C:

	}
	return x
}

/**
  time.After oom Verify demo
*/
func main() {
	ch := make(chan string, 100)

	go func() {
		for {
			ch <- "asong"
		}
	}()
	go func() {
		// Open pprof to listen for requests
		ip := "127.0.0.1:6060"
		if err := http.ListenAndServe(ip, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", ip)
		}
	}()

	for {
		add(ch)
	}
}
