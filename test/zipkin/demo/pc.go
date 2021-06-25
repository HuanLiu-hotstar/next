package main

import (
	// "context"
	"fmt"
	"io/ioutil"
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
	serverName = "PC"
	localAddr  = "192.168.1.63"
	port       = 8083
	umaddr     = "192.168.1.61:61"
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
	mux.HandleFunc("/pc", pc)
	// spanName := "root"
	middler := zipkinmw.NewServerMiddleware(
		tracer,
		// zipkinmw.SpanName(spanName),
		zipkinmw.TagResponseSize(true),
	)

	h := middler(mux)
	log.Printf("Server listening! %d ...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), h))

}
func pc(w http.ResponseWriter, r *http.Request) {
	span, _ := tracer.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	doclient(r)
	x := rand.Intn(100) + 10
	w.Write([]byte(fmt.Sprintf(`{"pc":%d}`, x)))
	//call um
}
func doclient(r *http.Request) {
	// create global zipkin traced http client
	co := zipkinmw.WithClient(&http.Client{Timeout: time.Second * 3})
	client, err := zipkinmw.NewClient(tracer, co, zipkinmw.ClientTrace(true), zipkinmw.ClientTags(map[string]string{"type:": "from-raw-http-client"}))
	if err != nil {
		log.Printf("unable to create client: %+v\n", err)
		return
	}

	// initiate a call to some_func
	url := fmt.Sprintf("http://%s/%s", umaddr, "/um")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("unable to create http request: %+v\n", err)
		return
	}

	req = req.WithContext(req.Context())
	res, err := client.DoWithAppSpan(req, "um")
	if err != nil {
		log.Printf("unable to do http request: %+v\n", err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	log.Printf("%s", body)
}
