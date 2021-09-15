package main

import (
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "time"

    "github.com/fsnotify/fsnotify"
    "gopkg.in/yaml.v2"
)

type Conf struct {
    Name string  `yaml:"name"`
    List []int32 `yaml:"list"`
}

func GetConf(filename string, conf *Conf) {
    f, err := os.Open(filename)
    if err != nil {
        log.Fatal("err:", err)
    }
    defer f.Close()
    bye, err := ioutil.ReadAll(f)
    if err != nil {
        log.Fatal(err)
    }

    err = yaml.Unmarshal(bye, conf)
    if err != nil {
        log.Printf("err:%s", err)
    }
}

func main() {
    filename := "./conf.yaml"
    watcher, err := fsnotify.NewWatcher()
    err = watcher.Add(filepath.Dir(filename))
    if err != nil {
        log.Fatal(err)
    }
    done := make(chan struct{}, 1)
    c := Conf{}
    GetConf(filename, &c)
    go func() {
        for {
            select {
            case <-done:
                return
            case e := <-watcher.Events:
                update := false
                if e.Op == fsnotify.Write || e.Op == fsnotify.Create {
                    update = true
                }
                if update {
                    GetConf(filename, &c)
                }
            }
        }
    }()
    for {
        select {

        case <-time.After(1 * time.Second):
            log.Print(c)

        }

    }

}
