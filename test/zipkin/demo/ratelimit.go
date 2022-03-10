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
	port   = ":18082"
)

func main() {

	localEndpoint, err := openzipkin.NewEndpoint("ratelimit", "192.168.1.61:8082")
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}
	// addr := "http://localhost:6831"
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	// reporter := zipkinHTTP.NewReporter(addr)
	defer reporter.Close()

	ratio := 0.001
	seed := int64(1000)
	sample, err := openzipkin.NewBoundarySampler(ratio, seed)
	if err != nil {
		log.Fatal("err:%s", err)
	}
	tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint), openzipkin.WithSampler(sample))
	if err != nil {
		panic(fmt.Sprintf("err:%s", err))
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/ratelimit", ratelimit)
	// spanName := "root"
	middler := zipkinmw.NewServerMiddleware(
		tracer,
		zipkinmw.SpanName("ratelimit"),
		zipkinmw.TagResponseSize(true),
	)

	h := middler(mux)
	log.Printf("Server listening! %s ...", port)
	log.Fatal(http.ListenAndServe(port, h))

}
func ratelimit(w http.ResponseWriter, r *http.Request) {
	span, _ := tracer.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	x := rand.Intn(100) + 100
	time.Sleep(time.Duration(x) * time.Millisecond)
	w.Write([]byte("hello ratelimit client"))
}
