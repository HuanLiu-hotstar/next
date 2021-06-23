package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinmw "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	tracer *openzipkin.Tracer
)

func main() {

	localEndpoint, err := openzipkin.NewEndpoint("golangsvc", "192.168.1.61:8082")
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	//exporter := zipkin.NewExporter(reporter, localEndpoint)
	tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))
	if err != nil {
		panic(fmt.Sprintf("err:%s", err))
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", list)
	spanName := "root"
	middler := zipkinmw.NewServerMiddleware(
		tracer,
		zipkinmw.SpanName(spanName),
		zipkinmw.TagResponseSize(true),
	)

	h := middler(mux)
	port := ":8080"
	log.Printf("Server listening! %s ...", port)
	log.Fatal(http.ListenAndServe(port, h))

}

func list(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	span, _ := tracer.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	//time.Sleep(6 * time.Second)
	database(r)
	serviceb(r)
	res := strings.Repeat("o", rand.Intn(100)+1)
	time.Sleep(time.Duration(rand.Intn(100)+1) * time.Millisecond)
	w.Write([]byte("Hello, w" + res + "rld!"))
}

func database(r *http.Request) {
	span, _ := tracer.StartSpanFromContext(r.Context(), "database")
	defer span.Finish()
	cache(r)
	x := rand.Intn(4) + 100
	time.Sleep(time.Duration(x) * time.Millisecond)
	span.Tag("sleep-time", fmt.Sprintf("database-cost:%d", x))
}

func cache(r *http.Request) {
	span, _ := tracer.StartSpanFromContext(r.Context(), "cache")
	defer span.Finish()
	x := rand.Intn(4) + 100
	time.Sleep(time.Duration(x) * time.Millisecond)
	span.Annotate(time.Now(), fmt.Sprintf("cost:%d", x))
}

func serviceb(r *http.Request) {
	span, pc := tracer.StartSpanFromContext(r.Context(), "serviceb")
	defer span.Finish()
	time.Sleep(time.Duration(rand.Intn(100)+100) * time.Millisecond)
	servicec(pc) // servicec is childof serviceb
	span.Annotate(time.Now(), "endtime")
}

//func servicec(r *http.Request) {
func servicec(c context.Context) {
	span, _ := tracer.StartSpanFromContext(c, "servicec")
	defer span.Finish()
	time.Sleep(time.Duration(rand.Intn(700)+100) * time.Millisecond)
	span.Tag("servicec", "C") // set tags for search servicec
}
