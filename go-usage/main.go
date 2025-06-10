package main

import (
	"encoding/json"
	"fmt"
)

type A struct {
	Landscape bool
	Portrait  bool
}
type B struct {
	A
	Landscape bool
}

func print(b B) {
	if b.Landscape {
		fmt.Println("b.landscape")
	}
	if b.A.Landscape {
		fmt.Println("b.A.landscape")
	}
	if b.Portrait {
		fmt.Println("b.A.Portrait")
	}
}
func main() {
	a := A{Landscape: true, Portrait: true}
	x := B{A: a, Landscape: false}
	bye, _ := json.Marshal(x)
	print(x)
	fmt.Printf("%s\n", bye)
	y := B{}
	json.Unmarshal(bye, &y)
	print(y)
}
