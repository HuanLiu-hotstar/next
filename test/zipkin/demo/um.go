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
	tracer     *openzipkin.Tracer
	serverName = "UM-Server"
	localAddr  = "192.168.1.61"
	port       = 8084
)

func main() {

	localEndpoint, err := openzipkin.NewEndpoint(serverName, fmt.Sprintf("%s:%d", localAddr, port))
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	tracer, err = openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(localEndpoint))
	if err != nil {
		panic(fmt.Sprintf("err:%s", err))
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/um", um)
	// spanName := "root"
	middler := zipkinmw.NewServerMiddleware(
		tracer,
		zipkinmw.SpanName(serverName),
		zipkinmw.TagResponseSize(true),
	)

	h := middler(mux)
	log.Printf("Server listening! %d ...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), h))

}
func um(w http.ResponseWriter, r *http.Request) {
	span, _ := tracer.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	x := rand.Intn(100) + 10
	time.Sleep(time.Duration(x) * time.Millisecond)
	w.Write([]byte(fmt.Sprintf(`{"um":%d}`, x)))
}
