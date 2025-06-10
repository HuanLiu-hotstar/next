package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Data struct {
	Key   string
	Value string
}

func (d *Data) Bytes() ([]byte, error) {
	//buf := bytes.Buffer{}
	bye, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return bye, nil
}
func read(file *os.File, offset int64) {

	bytesLen := make([]byte, 4)
	n, err := file.ReadAt(bytesLen, int64(offset))
	if err != nil {
		log.Fatal("err:%s", err)
	}
	log.Printf("len bytes:%d ", n)

	len := 0
	bytesBuf := bytes.NewBuffer(bytesLen)
	binary.Read(bytesBuf, binary.BigEndian, &len)
	log.Printf("len:%d", len)
	data := make([]byte, len)
	file.ReadAt(data, offset+int64(len))

}

func write(file *os.File, data []byte) {

	bytesBuf := bytes.NewBuffer([]byte{})
	len := len(data)
	binary.Write(bytesBuf, binary.BigEndian, len)
	n, err := bytesBuf.Write(data)
	if err != nil {
		log.Fatal("err:%s", err)
	}
	log.Printf("len:%d %d", n, n+4)
	n, err = file.Write(bytesBuf.Bytes())
	if err != nil {
		log.Fatal("err:%s", err)
	}
	log.Printf("size:%d", n)
}

type ValueCtx struct {
	Key string
}

func main() {
	file, err := os.OpenFile("test.data", os.O_CREATE|os.O_RDONLY|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("err:%s", err)
	}

	defer file.Close()

	buf := []byte("123")
	write(file, buf)
	a, b := ValueCtx{Key: "hello"}, ValueCtx{Key: "hello"}
	fmt.Println(a == b)
	//	file.WriteAt(buf, off)
}
