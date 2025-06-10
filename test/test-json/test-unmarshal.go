package main

import (
	"encoding/json"
	"log"
)

func main() {
	var v interface{}
	b := `{
		"hello":1
	}`
	err := json.Unmarshal([]byte(b), v)
	if err != nil {
		log.Printf("err:%s", err)
	}
	log.Printf("%v", v)

}
