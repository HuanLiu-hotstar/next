package main

import (
	// "contrib.go.opencensus.io/exporter/prometheus"
	"contrib.go.opencensus.io/exporter/zipkin"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"

	zipkinmw "github.com/openzipkin/zipkin-go/middleware/http"

	// "go.opencensus.io/plugin/ochttp"
	// "go.opencensus.io/stats/view"
	"go.opencensus.io/trace"

	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func registerZipkin() {
	localEndpoint, err := openzipkin.NewEndpoint("golangsvc", "192.168.1.61:8082")
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	exporter := zipkin.NewExporter(reporter, localEndpoint)
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}
func main() {
	// registerZipkin()
	localEndpoint, err := openzipkin.NewEndpoint("golangsvc", "192.168.1.61:8082")
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	//exporter := zipkin.NewExporter(reporter, localEndpoint)
	tracer, err := openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))
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
	log.Printf("Server listening! %s...", port)
	log.Fatal(http.ListenAndServe(port, h))
}
func list(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	//time.Sleep(6 * time.Second)
	database(r)
	serviceb(r)
	res := strings.Repeat("o", rand.Intn(100)+1)
	time.Sleep(time.Duration(rand.Intn(100)+1) * time.Millisecond)
	w.Write([]byte("Hello, w" + res + "rld!"))
}
func database(r *http.Request) {
	cache(r)
	_, span := trace.StartSpan(r.Context(), "database")
	defer span.End()
	time.Sleep(time.Duration(rand.Intn(4)+100) * time.Millisecond)
}
func cache(r *http.Request) {
	_, span := trace.StartSpan(r.Context(), "cache")
	defer span.End()
	time.Sleep(time.Duration(rand.Intn(100)+1) * time.Millisecond)
}
func serviceb(r *http.Request) {
	_, span := trace.StartSpan(r.Context(), "serviceb")
	defer span.End()
	time.Sleep(time.Duration(rand.Intn(100)+100) * time.Millisecond)
	servicec(r)
}
func servicec(r *http.Request) {
	_, span := trace.StartSpan(r.Context(), "servicec")
	defer span.End()
	time.Sleep(time.Duration(rand.Intn(700)+100) * time.Millisecond)
}
