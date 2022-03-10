package main

import (
	// "context"
	"fmt"
	// "io/ioutil"
	"log"
	"math/rand"
	"net/http"
	// "strings"
	"time"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinmw "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	tracer *openzipkin.Tracer
	port   = ":18085"
)

func main() {

	localEndpoint, err := openzipkin.NewEndpoint("auth", "localhost:8085")
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}
	addr := "http://localhost:6831"
	// /reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	reporter := zipkinHTTP.NewReporter(addr)
	defer reporter.Close()

	tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))
	if err != nil {
		panic(fmt.Sprintf("err:%s", err))
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", auth)
	mux.HandleFunc("/auth2", auth2)
	// spanName := "root"
	middler := zipkinmw.NewServerMiddleware(
		tracer,
		zipkinmw.SpanName("auth"),
		zipkinmw.TagResponseSize(true),
	)

	h := middler(mux)
	log.Printf("Server listening! %s ...", port)
	log.Fatal(http.ListenAndServe(port, h))

}
func auth(w http.ResponseWriter, r *http.Request) {
	span, _ := tracer.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	log.Printf("path:%s", r.URL.Path)
	x := rand.Intn(100) + 30
	time.Sleep(time.Duration(x) * time.Millisecond)
	w.Write([]byte("hello auth client"))
}

func auth2(w http.ResponseWriter, r *http.Request) {
	span, _ := tracer.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	log.Printf("path:%s", r.URL.Path)
	x := rand.Intn(100) + 30
	time.Sleep(time.Duration(x) * time.Millisecond)
	w.Write([]byte("hello auth2 client"))
}
