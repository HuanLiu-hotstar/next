package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Data struct {
	Key   string
	Value string
}
type Record struct {
	data   Data
	size   int
	offset int
}
type Hash struct {
	hash   map[string]Record
	r      *os.File
	w      *bufio.Writer
	offset int
}

func NewHash(r *os.File) *Hash {
	return &Hash{
		hash:   make(map[string]Record),
		r:      r,
		w:      bufio.NewWriter(r),
		offset: 0,
	}
}

func (h *Hash) Add(key, value string) error {
	r := Record{
		data:   Data{key, value},
		offset: h.offset,
	}
	bye, err := json.Marshal(r.data)
	if err != nil {
		return err
	}
	n, err := h.w.Write(bye)
	if err != nil {
		return err
	}
	r.size = len(bye)
	h.offset += n
	h.hash[key] = r
	log.Printf("success add %s %s %v", key, value, r)
	return nil
}

func (h *Hash) Get(key string) (string, error) {
	if _, ok := h.hash[key]; !ok {
		return "", nil
	}
	r := h.hash[key]
	p := make([]byte, r.size)
	log.Printf("r:%v", r)
	n, err := h.r.ReadAt(p, int64(r.offset))
	if err != nil || n != len(p) {
		return "", fmt.Errorf("err:%s len:%d", err, n)
	}

	d := Data{}
	err = json.Unmarshal(p, &d)
	if err != nil {
		return "", err
	}
	return d.Value, nil
}

func (h *Hash) Flush() {
	h.w.Flush()
}

func main() {
	// f, err := os.OpenFile("hashdata.log", os.O_CREATE|os.O_APPEND, 0666)
	f, err := os.OpenFile("hashdata.log", os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	// f.WriteString("hello world")
	defer f.Close()
	hash := NewHash(f)
	// N := 10
	// for i := 0; i < N; i++ {
	// 	key := fmt.Sprintf("k-%d", i)
	// 	value := fmt.Sprintf("v-%d", i+1)
	// 	hash.Add(key, value)
	// }
	// hash.Flush()
	key := fmt.Sprintf("k-%d", 0)
	value, err := hash.Get(key)
	log.Printf("value:%s,err:%s", value, err)
	// for i := 0; i < N; i++ {
	// 	key := fmt.Sprintf("k-%d", i)

	// 	value, err := hash.Get(key)
	// 	if err != nil {
	// 		log.Printf("err read key:%s %s", key, err)
	// 		continue
	// 	}
	// 	log.Printf("%s, v is %s", key, value)
	// }

}
