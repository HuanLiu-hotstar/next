package main

import (
	"encoding/json"
	"log"
	// "strings"
)
type Data struct {
	X int `json:"x"`
}
func main() {
	v := new(Data)
	b := "null"
	err := json.Unmarshal([]byte(b), v) // &v
	if err != nil {
		log.Printf("err:%s", err)
	}
	log.Printf("%v", v)
	y := v.X
	log.Printf("%d",y)
	// err = json.NewDecoder(strings.NewReader(b)).Decode(&v)
	// if err != nil {
	// 	log.Printf("err:%s",err)
	// }
	// log.Printf("v2:%d",v.X)


	var w *Data 
	bye, _ :=json.Marshal(w)
	log.Println(string(bye))

}
