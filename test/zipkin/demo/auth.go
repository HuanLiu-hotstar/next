package main

import (
	"context"
	"fmt"
	"io/ioutil"
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

	localEndpoint, err := openzipkin.NewEndpoint("pc", "192.168.1.61:8082")
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
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
	mux.HandleFunc("/pc", pc)
	// spanName := "root"
	middler := zipkinmw.NewServerMiddleware(
		tracer,
		// zipkinmw.SpanName(spanName),
		zipkinmw.TagResponseSize(true),
	)

	h := middler(mux)
	port := ":8080"
	log.Printf("Server listening! %s ...", port)
	log.Fatal(http.ListenAndServe(port, h))

}
func pc(w http.ResponseWriter, r *http.Request) {
	span, _ := tracer.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	//call auth
	//call rate limit
	//call pc
	//w.Write([]byte("hello client"))
}
