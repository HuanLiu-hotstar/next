package main

import "fmt"

type Data struct {
	Level int
	Name  string
	Z     *Other
}

type Other struct {
	ID int
}

func main() {
	d := &Data{Level: 1, Name: "test"}
	s := fmt.Sprintf("%+v", d)
	fmt.Println(s)
	y := []*Data{d, {Z: &Other{ID: 1}}}
	fmt.Printf("%+v", y[1].Z)
}
